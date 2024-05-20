package settings

import "context"

type settingsKey struct{}

func WithSettings(ctx context.Context, set Settings) context.Context {
	return context.WithValue(ctx, settingsKey{}, set)
}

func GetSettings(ctx context.Context) Settings {
	v, _ := ctx.Value(settingsKey{}).(Settings)
	return v
}
