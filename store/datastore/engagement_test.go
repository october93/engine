package datastore_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
	"github.com/october93/engine/store/datastore"
	"github.com/october93/engine/test"
)

func TestEngagements(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_engagement_test"
	db := test.DBInit(t, cfg)
	if os.Getenv("DEBUG") != "1" {
		// drop the database after the test is finished
		defer test.DBCleanup(t, db)
	}

	store, err := datastore.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Close the store
	defer func() {
		if e := store.Close(); e != nil {
			t.Fatalf("store.Close() failed.: %s", e)
		}
	}()

	chad := &model.User{
		ID:          "fabff660-1e07-4596-9247-f5c0656d5d36",
		Email:       "chad@october.news",
		Username:    "chad",
		DisplayName: "Chad Unicorn",
	}
	richard := &model.User{
		ID:               "0b0ce4ed-abed-4357-accb-8909036d01b2",
		Email:            "richard@piedpiper.com",
		Username:         "richard",
		DisplayName:      "Richard Hendricks",
		ProfileImagePath: "richard.png",
	}
	erlich := &model.User{
		ID:               "52df8f20-8b15-4d9c-a3c7-2d5417df2346",
		Email:            "erlich@piedpiper.com",
		Username:         "erlich",
		DisplayName:      "Erlich Bachman",
		ProfileImagePath: "erlich.png",
	}
	dinesh := &model.User{
		ID:               "179d1823-f025-4222-8ccc-bc18ba98cdf5",
		Email:            "dinesh@piedpiper.com",
		Username:         "dinesh",
		DisplayName:      "Dinesh Chugtai",
		ProfileImagePath: "dinesh.png",
	}
	betram := &model.User{
		ID:               "442bd9a0-73f5-4233-85a1-5aee260c2623",
		Email:            "bertram@piedpiper.com",
		Username:         "betram",
		DisplayName:      "Betram Gilfoyle",
		ProfileImagePath: "gilfoyle.png",
	}
	monica := &model.User{
		ID:               "745e6578-50c9-4a8e-b110-6ea4a52f0cc2",
		Email:            "monica@piedpiper.com",
		Username:         "monica",
		DisplayName:      "Monica Hall",
		ProfileImagePath: "monica.png",
	}
	gavin := &model.User{
		ID:               "a099bdd9-43e0-426c-a783-6452e210fcdc",
		Email:            "gavin.belson@hooli.com",
		Username:         "Gavin",
		DisplayName:      "Gavin Belson",
		ProfileImagePath: "gavin.jpg",
	}
	nelson := &model.User{
		ID:               "4b6b12c7-08a1-40e9-b908-5b7a0acc533b",
		Email:            "nelson@hooli.com",
		Username:         "bighead",
		DisplayName:      "Nelson 'Big Head' Bighetti",
		ProfileImagePath: "nelson.bmp",
	}

	alias1 := &model.AnonymousAlias{
		ID:               "40c5a86a-ab1a-4174-9096-afdc40dd862b",
		Username:         "egg",
		DisplayName:      "Anonymous",
		ProfileImagePath: "egg.png",
	}
	alias2 := &model.AnonymousAlias{
		ID:               "5858f233-b6e1-40f5-a0ac-a73be339f1d6",
		Username:         "mouse",
		DisplayName:      "Anonymous",
		ProfileImagePath: "mouse.png",
	}

	now := time.Now().UTC()

	card1 := &model.Card{
		ID:        "382319f1-778e-4eca-b48d-a25bce61f3f8",
		OwnerID:   chad.ID,
		CreatedAt: now.Add(1 * time.Minute),
	}
	card2 := &model.Card{
		ID:        "fa082e10-ebcf-4c52-ba7a-d36966f2b883",
		OwnerID:   richard.ID,
		CreatedAt: now.Add(3 * time.Minute),
	}
	card3 := &model.Card{
		ID:        "2956720a-0ecf-4d69-a853-c8a8c32080fc",
		OwnerID:   erlich.ID,
		CreatedAt: now.Add(4 * time.Minute),
	}
	card4 := &model.Card{
		ID:        "32ae4de4-63f9-4b47-b34a-fa6035361929",
		OwnerID:   betram.ID,
		CreatedAt: now.Add(5 * time.Minute),
	}
	card5 := &model.Card{
		ID:        "0d85b2f6-3e9e-47dc-8e2b-ab7a742cbed5",
		AliasID:   alias2.ID,
		OwnerID:   monica.ID,
		CreatedAt: now.Add(8 * time.Minute),
	}
	card6 := &model.Card{
		ID:        "b49af5c9-5880-4663-ae89-2882924cf53b",
		OwnerID:   monica.ID,
		CreatedAt: now.Add(9 * time.Minute),
	}
	card7 := &model.Card{
		ID:        "75ec8650-c2f5-4d15-aa1c-21e76d142104",
		OwnerID:   gavin.ID,
		CreatedAt: now.Add(10 * time.Minute),
		DeletedAt: model.NewDBTime(now.Add(11 * time.Minute)),
	}

	card2.ReplyTo(card1)
	card3.ReplyTo(card2)
	card5.ReplyTo(card3)
	card6.ReplyTo(card5)
	card7.ReplyTo(card1)

	reaction1 := &model.UserReaction{
		UserID:    richard.ID,
		CardID:    card1.ID,
		Type:      model.ReactionLike,
		CreatedAt: now.Add(2 * time.Minute),
	}
	reaction2 := &model.UserReaction{
		UserID:    betram.ID,
		CardID:    card1.ID,
		Type:      model.ReactionLike,
		CreatedAt: now.Add(1 * time.Minute),
	}
	reaction3 := &model.UserReaction{
		UserID:    card1.OwnerID,
		CardID:    card1.ID,
		Type:      model.ReactionLike,
		CreatedAt: now.Add(6 * time.Minute),
	}
	reaction4 := &model.UserReaction{
		UserID:    dinesh.ID,
		AliasID:   alias1.ID,
		CardID:    card1.ID,
		Type:      model.ReactionLike,
		CreatedAt: now.Add(7 * time.Minute),
	}

	setupUsers(t, store, chad, richard, erlich, dinesh, betram, monica, gavin, nelson)
	setupAnonymousAliases(t, store, alias1, alias2)
	setupCards(t, store, card1, card2, card3, card4, card5, card6, card7)
	setupReactions(t, store, reaction1, reaction2, reaction3, reaction4)
	expected := &model.Engagement{
		EngagedUsersByType: model.EngagedUsersByType{
			Comment: []*model.EngagedUser{
				&model.EngagedUser{
					UserID:           monica.ID,
					ProfileImagePath: monica.ProfileImagePath,
					Username:         monica.Username,
					DisplayName:      monica.DisplayName,
					Type:             "comment",
					CreatedAt:        now.Add(9 * time.Minute).Round(time.Microsecond),
				},
				&model.EngagedUser{
					UserID:           erlich.ID,
					ProfileImagePath: erlich.ProfileImagePath,
					Username:         erlich.Username,
					DisplayName:      erlich.DisplayName,
					Type:             "comment",
					CreatedAt:        now.Add(4 * time.Minute).Round(time.Microsecond),
				},
				&model.EngagedUser{
					UserID:           richard.ID,
					ProfileImagePath: richard.ProfileImagePath,
					Username:         richard.Username,
					DisplayName:      richard.DisplayName,
					Type:             "comment",
					CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
				},
				&model.EngagedUser{
					AliasID:          alias2.ID,
					ProfileImagePath: alias2.ProfileImagePath,
					Username:         alias2.Username,
					DisplayName:      alias2.DisplayName,
					Type:             "comment",
					CreatedAt:        now.Add(8 * time.Minute).Round(time.Microsecond),
				},
			},
			Like: []*model.EngagedUser{
				&model.EngagedUser{
					UserID:           chad.ID,
					ProfileImagePath: chad.ProfileImagePath,
					Username:         chad.Username,
					DisplayName:      chad.DisplayName,
					Type:             "like",
					CreatedAt:        now.Add(6 * time.Minute).Round(time.Microsecond),
				},
				&model.EngagedUser{
					UserID:           richard.ID,
					ProfileImagePath: richard.ProfileImagePath,
					Username:         richard.Username,
					DisplayName:      richard.DisplayName,
					Type:             "like",
					CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
				},
				&model.EngagedUser{
					UserID:           betram.ID,
					ProfileImagePath: betram.ProfileImagePath,
					Username:         betram.Username,
					DisplayName:      betram.DisplayName,
					Type:             "like",
					CreatedAt:        now.Add(1 * time.Minute).Round(time.Microsecond),
				},
				&model.EngagedUser{
					AliasID:          alias1.ID,
					ProfileImagePath: alias1.ProfileImagePath,
					Username:         alias1.Username,
					DisplayName:      alias1.DisplayName,
					Type:             "like",
					CreatedAt:        now.Add(7 * time.Minute).Round(time.Microsecond),
				},
			},
		},
		Count:        8,
		CommentCount: 4,
	}

	t.Run("get engagement", func(t *testing.T) {
		var engagement *model.Engagement
		engagement, err = store.GetEngagement(card1.ID) // nelson.ID
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(engagement, expected) {
			t.Errorf("unexpected result, diff: %v", pretty.Diff(engagement, expected))
		}
	})

	t.Run("no engagement", func(t *testing.T) {
		var engagement *model.Engagement
		engagement, err = store.GetEngagement(card7.ID) // nelson.ID
		if err != nil {
			t.Fatal(err)
		}
		if err != nil {
			t.Fatal(err)
		}
		if engagement != nil {
			t.Errorf("expected no engagement, actual %v", engagement)
		}
	})
}

func TestGetEngagement(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping - short test detected")
	}

	chad := &model.User{
		ID:          "fabff660-1e07-4596-9247-f5c0656d5d36",
		Email:       "chad@october.news",
		Username:    "chad",
		DisplayName: "Chad Unicorn",
	}
	richard := &model.User{
		ID:               "0b0ce4ed-abed-4357-accb-8909036d01b2",
		Email:            "richard@piedpiper.com",
		Username:         "richard",
		DisplayName:      "Richard Hendricks",
		ProfileImagePath: "richard.png",
	}
	erlich := &model.User{
		ID:               "52df8f20-8b15-4d9c-a3c7-2d5417df2346",
		Email:            "erlich@piedpiper.com",
		Username:         "erlich",
		DisplayName:      "Erlich Bachman",
		ProfileImagePath: "erlich.png",
	}
	dinesh := &model.User{
		ID:               "179d1823-f025-4222-8ccc-bc18ba98cdf5",
		Email:            "dinesh@piedpiper.com",
		Username:         "dinesh",
		DisplayName:      "Dinesh Chugtai",
		ProfileImagePath: "dinesh.png",
	}
	betram := &model.User{
		ID:               "442bd9a0-73f5-4233-85a1-5aee260c2623",
		Email:            "bertram@piedpiper.com",
		Username:         "betram",
		DisplayName:      "Betram Gilfoyle",
		ProfileImagePath: "gilfoyle.png",
	}
	monica := &model.User{
		ID:               "745e6578-50c9-4a8e-b110-6ea4a52f0cc2",
		Email:            "monica@piedpiper.com",
		Username:         "monica",
		DisplayName:      "Monica Hall",
		ProfileImagePath: "monica.png",
	}
	gavin := &model.User{
		ID:               "a099bdd9-43e0-426c-a783-6452e210fcdc",
		Email:            "gavin.belson@hooli.com",
		Username:         "Gavin",
		DisplayName:      "Gavin Belson",
		ProfileImagePath: "gavin.jpg",
	}
	nelson := &model.User{
		ID:               "4b6b12c7-08a1-40e9-b908-5b7a0acc533b",
		Email:            "nelson@hooli.com",
		Username:         "bighead",
		DisplayName:      "Nelson 'Big Head' Bighetti",
		ProfileImagePath: "nelson.bmp",
	}

	egg := &model.AnonymousAlias{
		ID:               "40c5a86a-ab1a-4174-9096-afdc40dd862b",
		Username:         "egg",
		DisplayName:      "Anonymous",
		ProfileImagePath: "egg.png",
	}
	mouse := &model.AnonymousAlias{
		ID:               "5858f233-b6e1-40f5-a0ac-a73be339f1d6",
		Username:         "mouse",
		DisplayName:      "Anonymous",
		ProfileImagePath: "mouse.png",
	}

	now := time.Now().UTC()

	tests := []struct {
		name      string
		root      *model.Card
		cardID    globalid.ID
		comments  []*model.Card
		reactions []*model.UserReaction
		expected  *model.Engagement
	}{
		{
			name: "no engagement",
			root: &model.Card{
				ID:        "6207277c-1318-4ffa-8885-5606544386c1",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			cardID:   "6207277c-1318-4ffa-8885-5606544386c1",
			expected: nil,
		},
		{
			name: "like public",
			root: &model.Card{
				ID:        "d6b9a519-81cc-4417-8d7c-52c28623daa8",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			cardID: "d6b9a519-81cc-4417-8d7c-52c28623daa8",
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "d6b9a519-81cc-4417-8d7c-52c28623daa8",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count: 1,
			},
		},
		{
			name: "3 likes public",
			root: &model.Card{
				ID:        "180e26f5-e1c3-4c21-aaae-6a5147c67b7b",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			cardID: "180e26f5-e1c3-4c21-aaae-6a5147c67b7b",
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "180e26f5-e1c3-4c21-aaae-6a5147c67b7b",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(1 * time.Minute),
				},
				&model.UserReaction{
					UserID:    erlich.ID,
					CardID:    "180e26f5-e1c3-4c21-aaae-6a5147c67b7b",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
				&model.UserReaction{
					UserID:    dinesh.ID,
					CardID:    "180e26f5-e1c3-4c21-aaae-6a5147c67b7b",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           dinesh.ID,
							ProfileImagePath: dinesh.ProfileImagePath,
							Username:         dinesh.Username,
							DisplayName:      dinesh.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(1 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count: 3,
			},
		},
		{
			name: "like anonymous",
			root: &model.Card{
				ID:        "265b7a95-2362-4a28-a1dc-0dee15f1fb24",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			cardID: "265b7a95-2362-4a28-a1dc-0dee15f1fb24",
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					AliasID:   egg.ID,
					CardID:    "265b7a95-2362-4a28-a1dc-0dee15f1fb24",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count: 1,
			},
		},
		{
			name:   "comment public",
			cardID: "ae3722dd-d7f8-443b-ab51-8930d623c6d5",
			root: &model.Card{
				ID:        "ae3722dd-d7f8-443b-ab51-8930d623c6d5",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "89e7cfef-5896-4b9f-9e6e-3c2666b11812",
					OwnerID:       richard.ID,
					ThreadRootID:  "ae3722dd-d7f8-443b-ab51-8930d623c6d5",
					ThreadReplyID: "ae3722dd-d7f8-443b-ab51-8930d623c6d5",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "comment anonymous",
			cardID: "030e4f27-76a7-448d-9dbc-5f568886e9b1",
			root: &model.Card{
				ID:        "030e4f27-76a7-448d-9dbc-5f568886e9b1",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "052cfb62-c9cd-4cd7-9618-ea69f3f20a0f",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "030e4f27-76a7-448d-9dbc-5f568886e9b1",
					ThreadReplyID: "030e4f27-76a7-448d-9dbc-5f568886e9b1",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "like public own post",
			cardID: "0203b1d7-f95f-443a-9556-a39225881bd7",
			root: &model.Card{
				ID:        "0203b1d7-f95f-443a-9556-a39225881bd7",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    chad.ID,
					CardID:    "0203b1d7-f95f-443a-9556-a39225881bd7",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           chad.ID,
							ProfileImagePath: chad.ProfileImagePath,
							Username:         chad.Username,
							DisplayName:      chad.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count: 1,
			},
		},
		{
			name:   "like anonymous own post",
			cardID: "8687e65b-e9c4-4295-9bf5-151e4edee476",
			root: &model.Card{
				ID:        "8687e65b-e9c4-4295-9bf5-151e4edee476",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    chad.ID,
					AliasID:   egg.ID,
					CardID:    "8687e65b-e9c4-4295-9bf5-151e4edee476",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count: 1,
			},
		},
		{
			name:   "comment public own post",
			cardID: "cdb6a18a-538c-4397-8bae-fa743efd03e1",
			root: &model.Card{
				ID:        "cdb6a18a-538c-4397-8bae-fa743efd03e1",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "de09abb0-8ee8-4b60-b9ea-d7787227abc5",
					OwnerID:       chad.ID,
					ThreadRootID:  "cdb6a18a-538c-4397-8bae-fa743efd03e1",
					ThreadReplyID: "cdb6a18a-538c-4397-8bae-fa743efd03e1",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           chad.ID,
							ProfileImagePath: chad.ProfileImagePath,
							Username:         chad.Username,
							DisplayName:      chad.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "comment anonymous own post",
			cardID: "14c5380f-2110-4a46-896e-10fbf23f70b4",
			root: &model.Card{
				ID:        "14c5380f-2110-4a46-896e-10fbf23f70b4",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "a65c57b9-e91c-400e-a478-4c94fbf27ef3",
					OwnerID:       chad.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "14c5380f-2110-4a46-896e-10fbf23f70b4",
					ThreadReplyID: "14c5380f-2110-4a46-896e-10fbf23f70b4",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "comment public, comment anonymous",
			cardID: "2ab5b392-8124-4e3f-adaa-7a2d8e2ae2af",
			root: &model.Card{
				ID:        "2ab5b392-8124-4e3f-adaa-7a2d8e2ae2af",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "de09abb0-8ee8-4b60-b9ea-d7787227abc5",
					OwnerID:       richard.ID,
					ThreadRootID:  "2ab5b392-8124-4e3f-adaa-7a2d8e2ae2af",
					ThreadReplyID: "2ab5b392-8124-4e3f-adaa-7a2d8e2ae2af",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "4343309e-7a03-46b2-9e5c-9be72c6f115a",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "2ab5b392-8124-4e3f-adaa-7a2d8e2ae2af",
					ThreadReplyID: "2ab5b392-8124-4e3f-adaa-7a2d8e2ae2af",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 2,
			},
		},
		{
			name:   "comment public, comment anonymous, comment public",
			cardID: "312bbe52-f32f-47d8-a3de-656f04882c39",
			root: &model.Card{
				ID:        "312bbe52-f32f-47d8-a3de-656f04882c39",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "4f59a39b-8322-4b85-93ea-0190b073e6b2",
					OwnerID:       richard.ID,
					ThreadRootID:  "312bbe52-f32f-47d8-a3de-656f04882c39",
					ThreadReplyID: "312bbe52-f32f-47d8-a3de-656f04882c39",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "12fbbc27-1d96-48b6-bea0-5b37265e3822",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "312bbe52-f32f-47d8-a3de-656f04882c39",
					ThreadReplyID: "312bbe52-f32f-47d8-a3de-656f04882c39",
					CreatedAt:     now.Add(3 * time.Minute),
				},
				&model.Card{
					ID:            "ac7bea47-92e6-4d5d-b015-fd127964c242",
					OwnerID:       richard.ID,
					ThreadRootID:  "312bbe52-f32f-47d8-a3de-656f04882c39",
					ThreadReplyID: "312bbe52-f32f-47d8-a3de-656f04882c39",
					CreatedAt:     now.Add(4 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 3,
			},
		},
		{
			name:   "like public, comment public",
			cardID: "9fcc7b6a-00ba-422f-8d96-d5ff11770e13",
			root: &model.Card{
				ID:        "9fcc7b6a-00ba-422f-8d96-d5ff11770e13",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "9fcc7b6a-00ba-422f-8d96-d5ff11770e13",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "0df7dc36-3b36-4607-868a-1583a3feac6d",
					OwnerID:       richard.ID,
					ThreadRootID:  "9fcc7b6a-00ba-422f-8d96-d5ff11770e13",
					ThreadReplyID: "9fcc7b6a-00ba-422f-8d96-d5ff11770e13",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "comment public, like public",
			cardID: "bc43a957-951e-48df-89bf-16e95b966968",
			root: &model.Card{
				ID:        "bc43a957-951e-48df-89bf-16e95b966968",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "33a4acac-a552-4076-b1cd-8eb5c6174a17",
					OwnerID:       richard.ID,
					ThreadRootID:  "bc43a957-951e-48df-89bf-16e95b966968",
					ThreadReplyID: "bc43a957-951e-48df-89bf-16e95b966968",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "bc43a957-951e-48df-89bf-16e95b966968",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "comment public, like public, like public",
			cardID: "f9c5b203-08ae-43d3-9631-edf642c70a51",
			root: &model.Card{
				ID:        "f9c5b203-08ae-43d3-9631-edf642c70a51",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "e960b9fe-9cdf-452d-a8c8-003fbe5b7259",
					OwnerID:       richard.ID,
					ThreadRootID:  "f9c5b203-08ae-43d3-9631-edf642c70a51",
					ThreadReplyID: "f9c5b203-08ae-43d3-9631-edf642c70a51",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "f9c5b203-08ae-43d3-9631-edf642c70a51",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
				&model.UserReaction{
					UserID:    erlich.ID,
					CardID:    "f9c5b203-08ae-43d3-9631-edf642c70a51",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(4 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(4 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        3,
				CommentCount: 1,
			},
		},
		{
			name:   "like anonymous, comment anonymous",
			cardID: "065b46a8-0708-4f8c-b986-d87db29fbfc0",
			root: &model.Card{
				ID:        "065b46a8-0708-4f8c-b986-d87db29fbfc0",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					AliasID:   egg.ID,
					CardID:    "065b46a8-0708-4f8c-b986-d87db29fbfc0",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "2627be32-bc08-424f-a696-ec9d2d7d02f3",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "065b46a8-0708-4f8c-b986-d87db29fbfc0",
					ThreadReplyID: "065b46a8-0708-4f8c-b986-d87db29fbfc0",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "comment anonymous, like anonymous",
			cardID: "76ed907c-63dd-4aec-99ca-160e3cefeba6",
			root: &model.Card{
				ID:        "76ed907c-63dd-4aec-99ca-160e3cefeba6",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "403a0da3-1044-40fe-ac23-d12eaf6744e2",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "76ed907c-63dd-4aec-99ca-160e3cefeba6",
					ThreadReplyID: "76ed907c-63dd-4aec-99ca-160e3cefeba6",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					AliasID:   egg.ID,
					CardID:    "76ed907c-63dd-4aec-99ca-160e3cefeba6",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},

		{
			name:   "like anonymous, comment public",
			cardID: "382319f1-778e-4eca-b48d-a25bce61f3f8",
			root: &model.Card{
				ID:        "382319f1-778e-4eca-b48d-a25bce61f3f8",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "89eecba1-21a3-447d-95c2-e98929e8d177",
					OwnerID:       richard.ID,
					ThreadRootID:  "382319f1-778e-4eca-b48d-a25bce61f3f8",
					ThreadReplyID: "382319f1-778e-4eca-b48d-a25bce61f3f8",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					AliasID:   egg.ID,
					CardID:    "382319f1-778e-4eca-b48d-a25bce61f3f8",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "like public, comment anonymous",
			cardID: "8c3bbc93-6d05-48c0-8b53-cc6ec5628ad4",
			root: &model.Card{
				ID:        "8c3bbc93-6d05-48c0-8b53-cc6ec5628ad4",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "331749ad-4f79-45fa-8c47-6850a1d25c63",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "8c3bbc93-6d05-48c0-8b53-cc6ec5628ad4",
					ThreadReplyID: "8c3bbc93-6d05-48c0-8b53-cc6ec5628ad4",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "8c3bbc93-6d05-48c0-8b53-cc6ec5628ad4",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{

						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "A likes public, B comments public",
			cardID: "32fb2ca2-1d43-4c40-a9d0-89303b8f3903",
			root: &model.Card{
				ID:        "32fb2ca2-1d43-4c40-a9d0-89303b8f3903",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "32fb2ca2-1d43-4c40-a9d0-89303b8f3903",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "0b1a4b75-c656-425c-81a9-cccc529e606c",
					OwnerID:       erlich.ID,
					ThreadRootID:  "32fb2ca2-1d43-4c40-a9d0-89303b8f3903",
					ThreadReplyID: "32fb2ca2-1d43-4c40-a9d0-89303b8f3903",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "A likes public, B comments public, C likes anonymous",
			cardID: "605fa1b6-df33-4c9c-89fb-7c14133d75f7",
			root: &model.Card{
				ID:        "605fa1b6-df33-4c9c-89fb-7c14133d75f7",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "605fa1b6-df33-4c9c-89fb-7c14133d75f7",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
				&model.UserReaction{
					UserID:    dinesh.ID,
					AliasID:   egg.ID,
					CardID:    "605fa1b6-df33-4c9c-89fb-7c14133d75f7",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(4 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "7c54ac5e-fe3f-4ace-8f88-a18a3e987748",
					OwnerID:       erlich.ID,
					ThreadRootID:  "605fa1b6-df33-4c9c-89fb-7c14133d75f7",
					ThreadReplyID: "605fa1b6-df33-4c9c-89fb-7c14133d75f7",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(4 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        3,
				CommentCount: 1,
			},
		},

		{
			name:   "one comment one deleted",
			cardID: "8d588b89-1c06-4adb-ba9c-755d66de691e",
			root: &model.Card{
				ID:        "8d588b89-1c06-4adb-ba9c-755d66de691e",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "6c7a65bd-99ad-4735-a531-30934c8d99aa",
					OwnerID:       richard.ID,
					ThreadRootID:  "8d588b89-1c06-4adb-ba9c-755d66de691e",
					ThreadReplyID: "8d588b89-1c06-4adb-ba9c-755d66de691e",
					CreatedAt:     now.Add(2 * time.Minute),
					DeletedAt:     model.NewDBTime(now.Add(3 * time.Minute)),
				},
			},
			expected: nil,
		},
		{
			name:   "two comments one deleted same author",
			cardID: "1aee21b3-5243-4164-99b8-729f6974b311",
			root: &model.Card{
				ID:        "1aee21b3-5243-4164-99b8-729f6974b311",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "ac42dc15-0bf7-4e26-a513-5635e36aadd5",
					OwnerID:       richard.ID,
					ThreadRootID:  "1aee21b3-5243-4164-99b8-729f6974b311",
					ThreadReplyID: "1aee21b3-5243-4164-99b8-729f6974b311",
					CreatedAt:     now.Add(1 * time.Minute),
				},
				&model.Card{
					ID:            "e4937f61-f66d-4aa7-93f7-bfa072153385",
					OwnerID:       richard.ID,
					ThreadRootID:  "1aee21b3-5243-4164-99b8-729f6974b311",
					ThreadReplyID: "1aee21b3-5243-4164-99b8-729f6974b311",
					CreatedAt:     now.Add(2 * time.Minute),
					DeletedAt:     model.NewDBTime(now.Add(3 * time.Minute)),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(1 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "two comments one deleted different authors",
			cardID: "c0073c65-6069-450f-80ed-5f22caa2d84c",
			root: &model.Card{
				ID:        "c0073c65-6069-450f-80ed-5f22caa2d84c",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "b8d9db63-e372-4ef7-95c3-ce08d939411a",
					OwnerID:       richard.ID,
					ThreadRootID:  "c0073c65-6069-450f-80ed-5f22caa2d84c",
					ThreadReplyID: "c0073c65-6069-450f-80ed-5f22caa2d84c",
					CreatedAt:     now.Add(1 * time.Minute),
				},
				&model.Card{
					ID:            "84cf1c7f-27c3-4d6a-9a5b-9c159e8ce628",
					OwnerID:       erlich.ID,
					ThreadRootID:  "c0073c65-6069-450f-80ed-5f22caa2d84c",
					ThreadReplyID: "c0073c65-6069-450f-80ed-5f22caa2d84c",
					CreatedAt:     now.Add(2 * time.Minute),
					DeletedAt:     model.NewDBTime(now.Add(3 * time.Minute)),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(1 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "public like, public comment (deleted)",
			cardID: "1af91253-b5d7-4e34-ad68-cda67756b058",
			root: &model.Card{
				ID:        "1af91253-b5d7-4e34-ad68-cda67756b058",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "1af91253-b5d7-4e34-ad68-cda67756b058",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "62b468b2-ba47-4010-a8ba-a1e63cf29711",
					OwnerID:       richard.ID,
					ThreadRootID:  "1af91253-b5d7-4e34-ad68-cda67756b058",
					ThreadReplyID: "1af91253-b5d7-4e34-ad68-cda67756b058",
					CreatedAt:     now.Add(3 * time.Minute),
					DeletedAt:     model.NewDBTime(now.Add(4 * time.Minute)),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 0,
			},
		},
		{
			name:   "anonymous like, public comment (deleted)",
			cardID: "d329c2ae-020b-40b4-87b9-871608aa6d01",
			root: &model.Card{
				ID:        "d329c2ae-020b-40b4-87b9-871608aa6d01",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					AliasID:   egg.ID,
					CardID:    "d329c2ae-020b-40b4-87b9-871608aa6d01",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "62b468b2-ba47-4010-a8ba-a1e63cf29711",
					OwnerID:       richard.ID,
					ThreadRootID:  "d329c2ae-020b-40b4-87b9-871608aa6d01",
					ThreadReplyID: "d329c2ae-020b-40b4-87b9-871608aa6d01",
					CreatedAt:     now.Add(3 * time.Minute),
					DeletedAt:     model.NewDBTime(now.Add(4 * time.Minute)),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 0,
			},
		},
		{
			name:   "public like, anonymous comment (deleted)",
			cardID: "b6414404-988f-4637-a14b-5bf0061ff237",
			root: &model.Card{
				ID:        "b6414404-988f-4637-a14b-5bf0061ff237",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    richard.ID,
					CardID:    "b6414404-988f-4637-a14b-5bf0061ff237",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(2 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "62b468b2-ba47-4010-a8ba-a1e63cf29711",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "b6414404-988f-4637-a14b-5bf0061ff237",
					ThreadReplyID: "b6414404-988f-4637-a14b-5bf0061ff237",
					CreatedAt:     now.Add(3 * time.Minute),
					DeletedAt:     model.NewDBTime(now.Add(4 * time.Minute)),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 0,
			},
		},
		{
			name:   "no engagement on comment",
			cardID: "f2f664ec-f6aa-452d-a212-b120a30f4d1f",
			root: &model.Card{
				ID:        "8469fb72-96bc-4561-99e3-c7171493262f",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "f2f664ec-f6aa-452d-a212-b120a30f4d1f",
					OwnerID:       richard.ID,
					ThreadRootID:  "8469fb72-96bc-4561-99e3-c7171493262f",
					ThreadReplyID: "8469fb72-96bc-4561-99e3-c7171493262f",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			expected: nil,
		},
		{
			name:   "public like on comment",
			cardID: "fa7390b0-adb2-41e8-9ea7-a40207e1d7d3",
			root: &model.Card{
				ID:        "626297f9-0d61-4e61-8a43-7bca56f29d54",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "fa7390b0-adb2-41e8-9ea7-a40207e1d7d3",
					OwnerID:       richard.ID,
					ThreadRootID:  "626297f9-0d61-4e61-8a43-7bca56f29d54",
					ThreadReplyID: "626297f9-0d61-4e61-8a43-7bca56f29d54",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    erlich.ID,
					CardID:    "fa7390b0-adb2-41e8-9ea7-a40207e1d7d3",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 0,
			},
		},
		{
			name:   "anonymous like on comment",
			cardID: "d81e7705-a2de-4fc1-9daf-e83cea2efd03",
			root: &model.Card{
				ID:        "c970287b-e359-4b5b-9af5-5cf88f11ce7b",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "d81e7705-a2de-4fc1-9daf-e83cea2efd03",
					OwnerID:       richard.ID,
					ThreadRootID:  "c970287b-e359-4b5b-9af5-5cf88f11ce7b",
					ThreadReplyID: "c970287b-e359-4b5b-9af5-5cf88f11ce7b",
					CreatedAt:     now.Add(2 * time.Minute),
				},
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    erlich.ID,
					AliasID:   egg.ID,
					CardID:    "d81e7705-a2de-4fc1-9daf-e83cea2efd03",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 0,
			},
		},
		{
			name:   "public comment on comment",
			cardID: "9a47cdd6-247c-4db9-b720-504fec189e2f",
			root: &model.Card{
				ID:        "093ee5d2-e004-46dd-ae5d-c8b05180dce6",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "9a47cdd6-247c-4db9-b720-504fec189e2f",
					OwnerID:       richard.ID,
					ThreadRootID:  "093ee5d2-e004-46dd-ae5d-c8b05180dce6",
					ThreadReplyID: "093ee5d2-e004-46dd-ae5d-c8b05180dce6",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "323d39bb-457f-4c74-89af-ba7c65f02d3f",
					OwnerID:       richard.ID,
					ThreadRootID:  "093ee5d2-e004-46dd-ae5d-c8b05180dce6",
					ThreadReplyID: "9a47cdd6-247c-4db9-b720-504fec189e2f",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "anonymous comment on comment",
			cardID: "897c05c5-9eac-483e-b695-426bce31b1e6",
			root: &model.Card{
				ID:        "21846458-82d1-4065-865a-73c8d822bf54",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "897c05c5-9eac-483e-b695-426bce31b1e6",
					OwnerID:       richard.ID,
					ThreadRootID:  "21846458-82d1-4065-865a-73c8d822bf54",
					ThreadReplyID: "21846458-82d1-4065-865a-73c8d822bf54",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "1fec4d46-9bd3-4e71-94d9-e619a5aa6242",
					OwnerID:       erlich.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "21846458-82d1-4065-865a-73c8d822bf54",
					ThreadReplyID: "897c05c5-9eac-483e-b695-426bce31b1e6",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        1,
				CommentCount: 1,
			},
		},
		{
			name:   "public comment and public like on comment",
			cardID: "d28c8dfd-5589-4c2c-a159-039d6680f8e0",
			root: &model.Card{
				ID:        "323d39bb-457f-4c74-89af-ba7c65f02d3f",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    erlich.ID,
					CardID:    "d28c8dfd-5589-4c2c-a159-039d6680f8e0",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "d28c8dfd-5589-4c2c-a159-039d6680f8e0",
					OwnerID:       richard.ID,
					ThreadRootID:  "323d39bb-457f-4c74-89af-ba7c65f02d3f",
					ThreadReplyID: "323d39bb-457f-4c74-89af-ba7c65f02d3f",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "b34d1c2d-6dc1-4e9c-87d5-66053433fc27",
					OwnerID:       richard.ID,
					ThreadRootID:  "323d39bb-457f-4c74-89af-ba7c65f02d3f",
					ThreadReplyID: "d28c8dfd-5589-4c2c-a159-039d6680f8e0",
					CreatedAt:     now.Add(4 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(4 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "public comment and anonymous like on comment",
			cardID: "ed474093-d8c3-419e-96e9-5aa455a4e217",
			root: &model.Card{
				ID:        "9e5a2128-6eb2-41c5-9686-7060a3d7d657",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    erlich.ID,
					AliasID:   egg.ID,
					CardID:    "ed474093-d8c3-419e-96e9-5aa455a4e217",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "ed474093-d8c3-419e-96e9-5aa455a4e217",
					OwnerID:       richard.ID,
					ThreadRootID:  "9e5a2128-6eb2-41c5-9686-7060a3d7d657",
					ThreadReplyID: "9e5a2128-6eb2-41c5-9686-7060a3d7d657",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "b34d1c2d-6dc1-4e9c-87d5-66053433fc27",
					OwnerID:       richard.ID,
					ThreadRootID:  "9e5a2128-6eb2-41c5-9686-7060a3d7d657",
					ThreadReplyID: "ed474093-d8c3-419e-96e9-5aa455a4e217",
					CreatedAt:     now.Add(4 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(4 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "anonymous comment and public like on comment",
			cardID: "cd93137a-0076-4062-af65-4da8ed416c10",
			root: &model.Card{
				ID:        "864dad93-c1ba-4bb0-8bac-4b92e9520c77",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			reactions: []*model.UserReaction{
				&model.UserReaction{
					UserID:    erlich.ID,
					CardID:    "cd93137a-0076-4062-af65-4da8ed416c10",
					Type:      model.ReactionLike,
					CreatedAt: now.Add(3 * time.Minute),
				},
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "cd93137a-0076-4062-af65-4da8ed416c10",
					OwnerID:       richard.ID,
					ThreadRootID:  "864dad93-c1ba-4bb0-8bac-4b92e9520c77",
					ThreadReplyID: "864dad93-c1ba-4bb0-8bac-4b92e9520c77",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "b34d1c2d-6dc1-4e9c-87d5-66053433fc27",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "864dad93-c1ba-4bb0-8bac-4b92e9520c77",
					ThreadReplyID: "cd93137a-0076-4062-af65-4da8ed416c10",
					CreatedAt:     now.Add(4 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(4 * time.Minute).Round(time.Microsecond),
						},
					},
					Like: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "like",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 1,
			},
		},
		{
			name:   "public comment and anonymous comment same author on comment",
			cardID: "7cd0f11d-f696-4cc1-a0b6-337e774c5d89",
			root: &model.Card{
				ID:        "1e0a6f31-1cd8-4f2b-b4dc-969c2783c4fe",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "7cd0f11d-f696-4cc1-a0b6-337e774c5d89",
					OwnerID:       richard.ID,
					ThreadRootID:  "1e0a6f31-1cd8-4f2b-b4dc-969c2783c4fe",
					ThreadReplyID: "1e0a6f31-1cd8-4f2b-b4dc-969c2783c4fe",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "6210aeb9-712e-4fca-8d4f-7f44831ba354",
					OwnerID:       richard.ID,
					ThreadRootID:  "1e0a6f31-1cd8-4f2b-b4dc-969c2783c4fe",
					ThreadReplyID: "7cd0f11d-f696-4cc1-a0b6-337e774c5d89",
					CreatedAt:     now.Add(3 * time.Minute),
				},
				&model.Card{
					ID:            "4630aa91-e34d-4a86-821a-7ee6c4c67c71",
					OwnerID:       richard.ID,
					AliasID:       egg.ID,
					ThreadRootID:  "1e0a6f31-1cd8-4f2b-b4dc-969c2783c4fe",
					ThreadReplyID: "7cd0f11d-f696-4cc1-a0b6-337e774c5d89",
					CreatedAt:     now.Add(4 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							AliasID:          egg.ID,
							ProfileImagePath: egg.ProfileImagePath,
							Username:         egg.Username,
							DisplayName:      egg.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(4 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 2,
			},
		},
		{
			name:   "comment on comment included in root engagement",
			cardID: "221b6fa0-2ca1-418c-8d30-6591bb9819b2",
			root: &model.Card{
				ID:        "221b6fa0-2ca1-418c-8d30-6591bb9819b2",
				OwnerID:   chad.ID,
				CreatedAt: now,
			},
			comments: []*model.Card{
				&model.Card{
					ID:            "95f1f90e-30f8-4e0f-b809-73d535cb06e8",
					OwnerID:       richard.ID,
					ThreadRootID:  "221b6fa0-2ca1-418c-8d30-6591bb9819b2",
					ThreadReplyID: "221b6fa0-2ca1-418c-8d30-6591bb9819b2",
					CreatedAt:     now.Add(2 * time.Minute),
				},
				&model.Card{
					ID:            "b798a36e-93c1-41e9-a110-b74e8bca45e8",
					OwnerID:       erlich.ID,
					ThreadRootID:  "221b6fa0-2ca1-418c-8d30-6591bb9819b2",
					ThreadReplyID: "95f1f90e-30f8-4e0f-b809-73d535cb06e8",
					CreatedAt:     now.Add(3 * time.Minute),
				},
			},
			expected: &model.Engagement{
				EngagedUsersByType: model.EngagedUsersByType{
					Comment: []*model.EngagedUser{
						&model.EngagedUser{
							UserID:           erlich.ID,
							ProfileImagePath: erlich.ProfileImagePath,
							Username:         erlich.Username,
							DisplayName:      erlich.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(3 * time.Minute).Round(time.Microsecond),
						},
						&model.EngagedUser{
							UserID:           richard.ID,
							ProfileImagePath: richard.ProfileImagePath,
							Username:         richard.Username,
							DisplayName:      richard.DisplayName,
							Type:             "comment",
							CreatedAt:        now.Add(2 * time.Minute).Round(time.Microsecond),
						},
					},
				},
				Count:        2,
				CommentCount: 2,
			},
		},
	}

	cfg := datastore.NewTestConfig()
	cfg.Database = "engine_datastore_get_engagement_test"
	db := test.DBInit(t, cfg)
	if os.Getenv("DEBUG") != "1" {
		// drop the database after the test is finished
		defer test.DBCleanup(t, db)
	}

	store, err := datastore.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Close the store
	defer func() {
		if e := store.Close(); e != nil {
			t.Fatalf("store.Close() failed.: %s", e)
		}
	}()

	setupUsers(t, store, chad, richard, erlich, dinesh, betram, monica, gavin, nelson)
	setupAnonymousAliases(t, store, egg, mouse)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupCards(t, store, tt.root)
			setupCards(t, store, tt.comments...)
			setupReactions(t, store, tt.reactions...)

			engagement, err := store.GetEngagement(tt.cardID) // tt.root.OwnerID
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(engagement, tt.expected) {
				t.Errorf("unexpected engagement, diff: %v", pretty.Diff(engagement, tt.expected))
			}
		})
	}
}

func setupUsers(t *testing.T, store *datastore.Store, users ...*model.User) {
	for _, user := range users {
		err := store.SaveUser(user)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupCards(t *testing.T, store *datastore.Store, cards ...*model.Card) {
	for _, card := range cards {
		err := store.SaveCard(card)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupReactions(t *testing.T, store *datastore.Store, reactions ...*model.UserReaction) {
	for _, reaction := range reactions {
		err := store.SaveUserReaction(reaction)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupAnonymousAliases(t *testing.T, store *datastore.Store, anonymousAliases ...*model.AnonymousAlias) {
	for _, anonymousAlias := range anonymousAliases {
		err := store.SaveAnonymousAlias(anonymousAlias)
		if err != nil {
			t.Fatal(err)
		}
	}
}
