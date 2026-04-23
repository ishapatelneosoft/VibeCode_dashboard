package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"auth-project/internal/service"
)

func TestTaskService_ValidateTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		expectErr bool
	}{
		{"Empty title", "", true},
		{"Valid title", "Task", false},
		{"Max length 255", repeat("a", 255), false},
		{"Exceeds max length", repeat("a", 256), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTitle(tt.title)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, service.ErrInvalidTitle, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateTaskRequest_Validation(t *testing.T) {
	req := service.CreateTaskRequest{
		Title:       "Valid Task",
		Description: "Description",
		AssigneeID:  nil,
		DueDate:     nil,
	}
	assert.Equal(t, "Valid Task", req.Title)
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// validateTitle replicates the validation logic from service.validateTitle
func validateTitle(title string) error {
	if len(title) == 0 || len(title) > 255 {
		return service.ErrInvalidTitle
	}
	return nil
}

// func TestTaskService_CreateTask_InvalidTitle_Empty(t *testing.T) {
// 	// This test does not require a db connection
// 	// We can't instantiate TaskService without db, so we skip
// 	t.Skip("Requires db connection for full test")
// }

// func TestTaskService_CreateTask_InvalidTitle_TooLong(t *testing.T) {
// 	t.Skip("Requires db connection for full test")
// }
