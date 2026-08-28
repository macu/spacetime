package spacetime

import (
	"database/sql"
	"fmt"
	"time"
)

const VoteWindow = 12 * time.Hour

func AddVote(db *sql.DB, userID uint, spaceID uint, voteValue int) (int, error) {
	// Add a vote to the database within the vote window

	if voteValue == 0 {
		return RemoveVote(db, userID, spaceID)
	}

	if voteValue < -1 || voteValue > 1 {
		return 0, fmt.Errorf("invalid vote value: %d", voteValue)
	}

	// Check if the user has already voted within the vote window
	var existingVoteId uint
	err := db.QueryRow(
		"SELECT id FROM space_vote WHERE user_id = $1 AND space_id = $2 AND created_at >= $3",
		userID, spaceID, time.Now().Add(-VoteWindow),
	).Scan(&existingVoteId)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("checking if vote exists: %w", err)
	}
	if existingVoteId != 0 {
		return 0, fmt.Errorf("user has already voted within the vote window")
	}

	// Insert the new vote into the database
	_, err = db.Exec(
		"INSERT INTO space_vote (user_id, space_id, vote_value, created_at) VALUES ($1, $2, $3, $4)",
		userID, spaceID, voteValue, time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add vote: %w", err)
	}
	return CalculateVoteSum(db, spaceID)
}

func RemoveVote(db *sql.DB, userID uint, spaceID uint) (int, error) {
	// Remove a vote from the database within the vote window

	var voteId uint
	err := db.QueryRow(
		"SELECT id FROM space_vote WHERE user_id = $1 AND space_id = $2 AND created_at >= $3",
		userID, spaceID, time.Now().Add(-VoteWindow),
	).Scan(&voteId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no vote found to remove")
		}
		return 0, fmt.Errorf("failed to find vote to remove: %w", err)
	}

	_, err = db.Exec("DELETE FROM space_vote WHERE id = $1", voteId)
	if err != nil {
		return 0, fmt.Errorf("failed to remove vote: %w", err)
	}

	return CalculateVoteSum(db, spaceID)

}

func CalculateVoteSum(db *sql.DB, spaceID uint) (int, error) {
	var voteSum int
	err := db.QueryRow(
		"SELECT COALESCE(SUM(vote_value), 0) FROM space_vote WHERE space_id = $1",
		spaceID,
	).Scan(&voteSum)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate vote sum: %w", err)
	}

	return voteSum, nil
}
