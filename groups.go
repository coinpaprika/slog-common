package slogcommon

import (
	"log/slog"
	"slices"

	"github.com/samber/lo"
)

func AppendAttrsToGroup(groups []string, actualAttrs []slog.Attr, newAttrs ...slog.Attr) []slog.Attr {
	actualAttrs = slices.Clone(actualAttrs)

	if len(groups) == 0 {
		return UniqAttrs(append(actualAttrs, newAttrs...))
	}

	for i := range actualAttrs {
		attr := actualAttrs[i]
		if attr.Key == groups[0] && attr.Value.Kind() == slog.KindGroup {
			actualAttrs[i] = slog.Group(groups[0], lo.ToAnySlice(AppendAttrsToGroup(groups[1:], attr.Value.Group(), newAttrs...))...)
			return actualAttrs
		}
	}

	return UniqAttrs(
		append(
			actualAttrs,
			slog.Group(
				groups[0],
				lo.ToAnySlice(AppendAttrsToGroup(groups[1:], []slog.Attr{}, newAttrs...))...,
			),
		),
	)
}

// @TODO: should be recursive
func UniqAttrs(attrs []slog.Attr) []slog.Attr {
	return slices.CompactFunc(attrs, func(a slog.Attr, b slog.Attr) bool { return a.Key == b.Key })
}
