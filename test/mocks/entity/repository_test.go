package mocks

import (
	"context"
	"fmt"
	entity "go-source/repositories/entity1"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestEntityRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIEntityRepository(ctrl)

	t.Run("Create_Entity", func(t *testing.T) {
		ctx := context.Background()
		mockData := &entity.Entity{
			Status: "123",
		}

		mockRepo.EXPECT().Create(gomock.Any(), mockData).Return(nil).Times(1)

		err := mockRepo.Create(ctx, mockData)
		assert.Equal(t, err, nil)
	})

	t.Run("Get_Entity", func(t *testing.T) {
		ctx := context.Background()
		objectID, err := primitive.ObjectIDFromHex("60c72b2f9af1c88f1f6b3b0c")
		if err != nil {
			fmt.Println("❌ Invalid ObjectID:", err)
			return
		}
		mockData := &entity.Entity{
			Id:     objectID,
			Status: "123",
		}

		mockRepo.EXPECT().Get(gomock.Any(), mockData.Id).Return(mockData, nil).Times(1)

		data, err := mockRepo.Get(ctx, objectID)
		assert.Equal(t, data, mockData)
		assert.Equal(t, err, nil)
	})
}
