package model_test

import (
	"ads/adservice/internal/domain/model"
	"ads/pkg/errs"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func vPtr[T any](v T) *T {
	return &v
}

func TestNewAd(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		sellerID    uuid.UUID
		title       string
		description *string
		price       int64
		expect      error
	}

	var (
		testSelID = uuid.New()
		testTitle = "Apartment in the center of Shanghai"
		testDesc  = "We are selling an apartment in the center of Shanghai."
		testPrice = int64(1000000)
	)

	var tests = []testCase{
		{
			name:        "success",
			sellerID:    testSelID,
			title:       testTitle,
			description: vPtr(testDesc),
			price:       testPrice,
			expect:      nil,
		},
		{
			name:     "nullable seller id",
			sellerID: uuid.Nil,
			expect:   errs.ErrValueIsInvalid,
		},
		{
			name:     "empty title",
			sellerID: testSelID,
			title:    "",
			expect:   errs.ErrValueIsRequired,
		},
		{
			name:     "invalid title",
			sellerID: testSelID,
			title:    "Sell", // a small string
			expect:   errs.ErrValueIsInvalid,
		},
		{
			name:        "empty description",
			sellerID:    testSelID,
			title:       testTitle,
			description: vPtr(""),
			expect:      errs.ErrValueIsRequired,
		},
		{
			name:        "invalid description",
			sellerID:    testSelID,
			title:       testTitle,
			description: vPtr(strings.Repeat(testDesc, 45)), // a large string
			expect:      errs.ErrValueIsInvalid,
		},
		{
			name:        "invalid price",
			sellerID:    testSelID,
			title:       testTitle,
			description: nil,
			price:       testPrice * -1, // negative price
			expect:      errs.ErrValueIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad, err := model.NewAd(
				tt.sellerID, tt.title,
				tt.description, tt.price,
			)
			if tt.expect == nil {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, ad.ID())
				assert.Equal(t, tt.sellerID, ad.SellerID())
				assert.Equal(t, tt.title, ad.Title())
				assert.Equal(t, tt.price, ad.Price())
				assert.Equal(t, ad.CreatedAt(), ad.UpdatedAt())
			} else {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, ad)
			}
		})
	}
}

func TestAd_Publish(t *testing.T) {
	t.Parallel()

	testAd := model.RestoreAd(
		uuid.New(), uuid.New(), "Sell a car", nil,
		int64(100000), model.AdOnModeration, time.Now(), time.Now(),
	)

	// Publish for the first time - correct
	err := testAd.Publish()
	require.NoError(t, err)
	require.Equal(t, model.AdPublished, testAd.Status())
	require.True(t, testAd.IsPublished())

	// Trying to publish again - failure
	err = testAd.Publish()
	require.Error(t, err)
}

func TestAd_Reject(t *testing.T) {
	t.Parallel()

	testAd := model.RestoreAd(
		uuid.New(), uuid.New(), "Sell a car", nil,
		int64(100000), model.AdOnModeration, time.Now(), time.Now(),
	)

	// Reject for the first time - correct
	err := testAd.Reject()
	require.NoError(t, err)
	require.Equal(t, model.AdRejected, testAd.Status())
	require.True(t, testAd.IsRejected())

	// Trying to reject again - failure
	err = testAd.Reject()
	require.Error(t, err)
}

func TestAd_Delete(t *testing.T) {
	t.Parallel()

	testAd := model.RestoreAd(
		uuid.New(), uuid.New(), "Sell a car", nil,
		int64(100000), model.AdPublished, time.Now(), time.Now(),
	)

	// Delete for the first time - correct
	err := testAd.Delete()
	require.NoError(t, err)
	require.Equal(t, model.AdDeleted, testAd.Status())
	require.True(t, testAd.IsDeleted())

	// Trying to delete again - failure
	err = testAd.Delete()
	require.Error(t, err)
}

func TestAd_Update(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		title       *string
		description *string
		price       *int64
		expect      error
	}

	var (
		testTitle = "Apartment in the center of Shanghai"
		testDesc  = "We are selling an apartment in the center of Shanghai."
		testPrice = int64(1000000)
	)

	var tests = []testCase{
		{
			name:        "success",
			title:       vPtr(testTitle),
			description: vPtr(testDesc),
			price:       vPtr(testPrice),
			expect:      nil,
		},
		{
			name:        "success - nothing to update",
			title:       nil,
			description: nil,
			price:       nil,
			expect:      nil,
		},
		{
			name:   "invalid title",
			title:  vPtr("Sell"), // a small string
			expect: errs.ErrValueIsInvalid,
		},
		{
			name:        "invalid description",
			title:       vPtr(testTitle),
			description: vPtr(strings.Repeat(testDesc, 45)), // a large string
			expect:      errs.ErrValueIsInvalid,
		},
		{
			name:        "invalid price",
			title:       vPtr(testTitle),
			description: nil,
			price:       vPtr(testPrice * -1), // negative price
			expect:      errs.ErrValueIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad, _ := model.NewAd(
				uuid.New(), "Shanghai night tour",
				vPtr("You will never forget it!"), int64(1000),
			)

			updAt := ad.UpdatedAt()
			time.Sleep(time.Millisecond) // wait to change time

			err := ad.Update(tt.title, tt.description, tt.price)

			if tt.expect == nil {
				require.NoError(t, err)
				if tt.title != nil {
					assert.Equal(t, *tt.title, ad.Title())
				}
				if tt.price != nil {
					assert.Equal(t, *tt.price, ad.Price())
				}
				if tt.description != nil {
					adDesc := ad.Description()
					assert.Equal(t, *tt.description, *adDesc)
				}
				assert.NotEqual(t, ad.UpdatedAt(), updAt)
			} else {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Equal(t, ad.UpdatedAt(), updAt)
			}
		})
	}
}
