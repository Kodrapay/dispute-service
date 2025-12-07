package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/kodra-pay/dispute-service/internal/models"
)

type DisputeRepository struct {
	db *sql.DB
}

func NewDisputeRepository(dsn string) (*DisputeRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &DisputeRepository{db: db}, nil
}

func (r *DisputeRepository) findTransactionIDByReference(ctx context.Context, ref string) (int, int, error) {
	var id, merchantID int
	err := r.db.QueryRowContext(ctx, `SELECT id, merchant_id FROM transactions WHERE reference = $1`, ref).Scan(&id, &merchantID)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	return id, merchantID, nil
}

func (r *DisputeRepository) Create(ctx context.Context, ref, reason string) (*models.Dispute, error) {
	txID, merchantID, err := r.findTransactionIDByReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	if txID == 0 {
		return nil, fmt.Errorf("transaction not found")
	}

	var dispute models.Dispute
	query := `
		INSERT INTO disputes (transaction_id, status, reason, evidence, opened_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, status, opened_at
	`
	evidenceJSON, _ := json.Marshal([]models.Evidence{})
	if err := r.db.QueryRowContext(ctx, query, txID, models.StatusOpen, reason, evidenceJSON).
		Scan(&dispute.ID, &dispute.Status, &dispute.OpenedAt); err != nil {
		return nil, err
	}
	dispute.TransactionID = txID
	dispute.Reference = ref
	dispute.MerchantID = merchantID
	return &dispute, nil
}

func (r *DisputeRepository) Get(ctx context.Context, id int) (*models.Dispute, error) {
	query := `
		SELECT d.id, d.transaction_id, d.status, d.reason, d.evidence, d.opened_at, d.closed_at, t.reference, t.merchant_id
		FROM disputes d
		JOIN transactions t ON d.transaction_id = t.id
		WHERE d.id = $1
	`
	var disp models.Dispute
	var evidenceJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&disp.ID,
		&disp.TransactionID,
		&disp.Status,
		&disp.Reason,
		&evidenceJSON,
		&disp.OpenedAt,
		&disp.ClosedAt,
		&disp.Reference,
		&disp.MerchantID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(evidenceJSON) > 0 {
		_ = json.Unmarshal(evidenceJSON, &disp.Evidence)
	}
	return &disp, nil
}

func (r *DisputeRepository) AddEvidence(ctx context.Context, id int, ev models.Evidence) (*models.Dispute, error) {
	disp, err := r.Get(ctx, id)
	if err != nil || disp == nil {
		return nil, err
	}
	disp.Evidence = append(disp.Evidence, ev)
	evidenceJSON, _ := json.Marshal(disp.Evidence)
	_, err = r.db.ExecContext(ctx, `UPDATE disputes SET evidence = $2 WHERE id = $1`, id, evidenceJSON)
	if err != nil {
		return nil, err
	}
	return disp, nil
}

func (r *DisputeRepository) ListByMerchant(ctx context.Context, merchantID int, limit int) ([]*models.Dispute, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT d.id, d.transaction_id, d.status, d.reason, d.evidence, d.opened_at, d.closed_at, t.reference, t.merchant_id
		FROM disputes d
		JOIN transactions t ON d.transaction_id = t.id
		WHERE t.merchant_id = $1
		ORDER BY d.opened_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Dispute
	for rows.Next() {
		var disp models.Dispute
		var evidenceJSON []byte
		if err := rows.Scan(
			&disp.ID, &disp.TransactionID, &disp.Status, &disp.Reason, &evidenceJSON,
			&disp.OpenedAt, &disp.ClosedAt, &disp.Reference, &disp.MerchantID,
		); err != nil {
			return nil, err
		}
		if len(evidenceJSON) > 0 {
			_ = json.Unmarshal(evidenceJSON, &disp.Evidence)
		}
		list = append(list, &disp)
	}
	return list, rows.Err()
}
