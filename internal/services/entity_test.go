package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	"go-source/internal/domains"
	entity "go-source/repositories/entity1"
	mocks "go-source/test/mocks/entity"
)

func TestEntityService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIEntityRepository(ctrl)
	service := NewEntityService(mockRepo)

	testCases := []struct {
		name        string
		inputID     string
		mockReturn  *entity.Entity
		mockError   error
		expected    *domains.Entity
		expectedErr bool
	}{
		{
			name:    "INVALID_PARAM",
			inputID: "60c72b2f9af1c88f1f6b3b0c",
			mockReturn: &entity.Entity{
				Id:     primitive.NewObjectID(), // This will be overridden inside test
				Status: "active",
			},
			mockError: nil,
			expected: &domains.Entity{
				Status: "active",
			},
			expectedErr: false,
		},
		{
			name:        "INVALID_OBJECT_ID",
			inputID:     "invalid-id",
			mockReturn:  nil,
			mockError:   nil,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "INVALID_ID",
			inputID:     "60c72b2f9af1c88f1f6b3b0d",
			mockReturn:  nil,
			mockError:   fmt.Errorf("repo error"),
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()

			objectID, err := primitive.ObjectIDFromHex(tc.inputID)
			if err == nil {
				if tc.mockReturn != nil {
					tc.mockReturn.Id = objectID
				}
				mockRepo.EXPECT().Get(ctx, objectID).Return(tc.mockReturn, tc.mockError)
			}

			data, err := service.Get(ctx, tc.inputID)
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Status, data.Status)
			}
		})
	}
}
