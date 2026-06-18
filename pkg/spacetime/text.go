package spacetime

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
)

func CreateText(tx *sql.Tx, auth ajax.Auth, parentID uint,
	text string, title *string, recording *NakedText, startTime *time.Time) (*Space, error) {

	if !ValidateText(text) {
		return nil, fmt.Errorf("invalid text")
	}

	if title != nil && !ValidateTitle(*title) {
		return nil, fmt.Errorf("invalid title")
	}

	if recording != nil && !ValidateNakedText(*recording, text) {
		return nil, fmt.Errorf("invalid text recording")
	}

	var space = &Space{
		ParentID:  &parentID,
		SpaceType: SpaceTypeText,
		Text:      &text,
		Title:     &title,
	}

	uniqueTextId, err := GetOrCreateUniqueTextId(tx, text)
	if err != nil {
		return nil, fmt.Errorf("get or create unique text id: %w", err)
	} else if uniqueTextId == nil {
		return nil, fmt.Errorf("unique text id is nil")
	}

	var uniqueTitleId *uint
	if title != nil {
		uniqueTitleIdValue, err := GetOrCreateUniqueTextId(tx, *title)
		if err != nil {
			return nil, fmt.Errorf("get or create unique text id for title: %w", err)
		} else if uniqueTitleIdValue == nil {
			return nil, fmt.Errorf("unique text id for title is nil")
		}
		uniqueTitleId = uniqueTitleIdValue
	}

	// Create space
	if err = CreateSpace(tx, auth, space, &parentID, SpaceTypeText); err != nil {
		return nil, fmt.Errorf("create space: %w", err)
	}

	var recordingJson *[]byte
	if recording != nil {
		recordingJsonData, err := json.Marshal(recording)
		if err != nil {
			return nil, fmt.Errorf("marshal text recording: %w", err)
		}
		recordingJson = &recordingJsonData
	} else {
		startTime = nil
	}

	// Create text_space
	if _, err = tx.Exec(`INSERT INTO text_space
			(space_id, parent_id, text_id, title_id, recording, started_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
		space.ID, parentID, *uniqueTextId, uniqueTitleId, recordingJson, startTime,
	); err != nil {
		return nil, fmt.Errorf("insert text_space: %w", err)
	}

	return space, nil

}

func LoadTextRecording(conn *sql.DB, space *Space) error {
	var jsonData []byte

	err := conn.QueryRow(`SELECT recording, started_at
		FROM text_space
		WHERE space_id = $1`, space.ID,
	).Scan(&jsonData, &space.StartedAt)

	if err != nil {
		return fmt.Errorf("load text recording: %w", err)
	}

	if jsonData != nil {
		var recording NakedText
		if err = json.Unmarshal(jsonData, &recording); err != nil {
			return fmt.Errorf("unmarshal text recording: %w", err)
		}
		space.ReplayData = &recording
		hasRecording := true
		space.HasRecording = &hasRecording
	} else {
		hasRecording := false
		space.HasRecording = &hasRecording
	}

	return nil
}
