package ui

//import (
//	"charm.land/lipgloss/v2"
//	"github.com/filz0r/jat/internal/config"
//	"github.com/google/uuid"
//)
//
//func centerOverlay(bg, overlay string, w, h int) string {
//	ow := lipgloss.Width(overlay)
//	oh := lipgloss.Height(overlay)
//	return lipgloss.NewCompositor(
//		lipgloss.NewLayer(bg).Z(0),
//		lipgloss.NewLayer(overlay).X(max(0, (w-ow)/2)).Y(max(0, (h-oh)/2)).Z(1),
//	).Render()
//}
//
//func parseConfigUserID(cfg *config.ConfigFile) (uuid.UUID, error) {
//	id, err := cfg.GetUserID()
//	if err != nil {
//		return uuid.UUID{}, err
//	}
//	return uuid.Parse(id)
//}
//
//// navigateList moves a selection index by delta within items, clamped to the
//// valid range. Shared by every tab so they all honour the same j/k/arrow
//// contract.
//func navigateList(items []string, selected, delta int) int {
//	return navigateIndex(len(items), selected, delta)
//}
//
//// navigateIndex is navigateList for callers that only have a count (modal
//// pickers, suggestion matches, scroll offsets).
//func navigateIndex(count, selected, delta int) int {
//	if count == 0 {
//		return 0
//	}
//	selected += delta
//	if selected < 0 {
//		selected = 0
//	}
//	if selected >= count {
//		selected = count - 1
//	}
//	return selected
//}
