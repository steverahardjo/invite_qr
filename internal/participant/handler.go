package participant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

func (h *Handler) ListParticipants(c *gin.Context) {
	rows, err := h.db.Query(c, "SELECT * FROM participants")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var participants []pgx.Row
	for rows.Next() {
		var p pgx.Row
		if err := rows.Scan(&p); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		participants = append(participants, p)
	}
	c.JSON(http.StatusOK, participants)
}
