package campaign

import "time"

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage // agar table campaign dapat berelasi dengan table campaign_images
}

type CampaignImage struct {
	ID               int
	CampaignID       int
	FileName         string
	IsPrimary        int
    CreatedAt        time.Time
	UpdatedAt        time.Time
}