package repository

import (
	"context"
	"database/sql"
)

type Repository interface {
	LogPoop(ctx context.Context, userID int64, username string, msgId int64, timestamp string, unixTimestamp int64) error
	GetGlobalPoopCount(ctx context.Context, userID int64) (int, error)
	GetMonthlyPoopCount(ctx context.Context, userID int64) (int, error)
	GetMonthlyPoopStats(ctx context.Context, userID int64) ([]MonthlyPoopCount, error)
	GetDaysWithoutPoop(ctx context.Context, userID int64) (int, error)
	GetMaxPoopStreak(ctx context.Context, userID int64) (int, error)
	GetDayWithMostPoops(ctx context.Context, userID int64) (string, int, error)
	GetMonthlyLeaderboard(ctx context.Context) ([]UserPoopCount, error)
	GetBottomPoopers(ctx context.Context) ([]UserPoopCount, error)
	GetMonthlyPoodium(ctx context.Context) ([]UserPoopCount, error)
	GetPastMonthPoodium(ctx context.Context) ([]UserPoopCount, error)
	GetYearlyPoodium(ctx context.Context) ([]UserPoopCount, error)
	GetYearlyPoopCount(ctx context.Context, userID int64, year int) (int, error)
	GetPoopsByHour(ctx context.Context, userID int64, year int) ([]HourDistribution, error)
	GetPoopsByDayOfWeek(ctx context.Context, userID int64, year int) ([]DayOfWeekDistribution, error)
	GetYearlyRanking(ctx context.Context, userID int64, year int) (YearlyRanking, error)
	GetGroupYearlyStats(ctx context.Context, year int) ([]UserPoopCount, error)
	GetGroupAwards(ctx context.Context, year int) ([]GroupAward, error)
	HealthCheck(ctx context.Context) error
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) Repository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) LogPoop(ctx context.Context, userID int64, username string, msgId int64, timestamp string, unixTimestamp int64) error {
	return LogPoop(ctx, r.db, userID, username, msgId, timestamp, unixTimestamp)
}

func (r *SQLiteRepository) GetGlobalPoopCount(ctx context.Context, userID int64) (int, error) {
	return GetGlobalPoopCount(ctx, r.db, userID)
}

func (r *SQLiteRepository) GetMonthlyPoopCount(ctx context.Context, userID int64) (int, error) {
	return GetMonthlyPoopCount(ctx, r.db, userID)
}

func (r *SQLiteRepository) GetMonthlyPoopStats(ctx context.Context, userID int64) ([]MonthlyPoopCount, error) {
	return GetMonthlyPoopStats(ctx, r.db, userID)
}

func (r *SQLiteRepository) GetDaysWithoutPoop(ctx context.Context, userID int64) (int, error) {
	return GetDaysWithoutPoop(ctx, r.db, userID)
}

func (r *SQLiteRepository) GetMaxPoopStreak(ctx context.Context, userID int64) (int, error) {
	return GetMaxPoopStreak(ctx, r.db, userID)
}

func (r *SQLiteRepository) GetDayWithMostPoops(ctx context.Context, userID int64) (string, int, error) {
	return GetDayWithMostPoops(ctx, r.db, userID)
}

func (r *SQLiteRepository) GetMonthlyLeaderboard(ctx context.Context) ([]UserPoopCount, error) {
	return GetMonthlyLeaderboard(ctx, r.db)
}

func (r *SQLiteRepository) GetBottomPoopers(ctx context.Context) ([]UserPoopCount, error) {
	return GetBottomPoopers(ctx, r.db)
}

func (r *SQLiteRepository) GetMonthlyPoodium(ctx context.Context) ([]UserPoopCount, error) {
	return GetMonthlyPoodium(ctx, r.db)
}

func (r *SQLiteRepository) GetPastMonthPoodium(ctx context.Context) ([]UserPoopCount, error) {
	return GetPastMonthPoodium(ctx, r.db)
}

func (r *SQLiteRepository) GetYearlyPoodium(ctx context.Context) ([]UserPoopCount, error) {
	return GetYearlyPoodium(ctx, r.db)
}

func (r *SQLiteRepository) GetYearlyPoopCount(ctx context.Context, userID int64, year int) (int, error) {
	return GetYearlyPoopCount(ctx, r.db, userID, year)
}

func (r *SQLiteRepository) GetPoopsByHour(ctx context.Context, userID int64, year int) ([]HourDistribution, error) {
	return GetPoopsByHour(ctx, r.db, userID, year)
}

func (r *SQLiteRepository) GetPoopsByDayOfWeek(ctx context.Context, userID int64, year int) ([]DayOfWeekDistribution, error) {
	return GetPoopsByDayOfWeek(ctx, r.db, userID, year)
}

func (r *SQLiteRepository) GetYearlyRanking(ctx context.Context, userID int64, year int) (YearlyRanking, error) {
	return GetYearlyRanking(ctx, r.db, userID, year)
}

func (r *SQLiteRepository) GetGroupYearlyStats(ctx context.Context, year int) ([]UserPoopCount, error) {
	return GetGroupYearlyStats(ctx, r.db, year)
}

func (r *SQLiteRepository) GetGroupAwards(ctx context.Context, year int) ([]GroupAward, error) {
	return GetGroupAwards(ctx, r.db, year)
}

func (r *SQLiteRepository) HealthCheck(ctx context.Context) error {
	return HealthCheck(ctx, r.db)
}
