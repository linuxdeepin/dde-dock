package desktop

import (
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/initializer"
	"pkg.deepin.io/lib/utils"
)

const (
	// FileManagerPerferenceSchemaID is filemanager's general preferences' schema id
	FileManagerPerferenceSchemaID string = "com.deepin.filemanager.preferences"

	// FileManagerDesktopSchemaID is desktop specific settings' schema id
	FileManagerDesktopSchemaID string = "com.deepin.dde.desktop"
)

const (
	// ConfirmTrash schema key.
	ConfirmTrash = "confirm-trash"

	// ConfirmEmptyTrash schema key.
	ConfirmEmptyTrash = "confirm-empty-trash"

	// ActivationPolicy schema key.
	ActivationPolicy = "activation-policy"

	// ClickPolicy schema key.
	ClickPolicy = "click-policy"

	// ShowThumbnail schema key.
	ShowThumbnail = "show-thumbnail"

	// ShowHiddenFiles schema key.
	ShowHiddenFiles = "show-hidden-files"

	// ShowExtensionName schema key.
	ShowExtensionName = "show-extension-name"

	// LabelPosition schema key.
	LabelPosition = "label-position"

	// AllowDeleteImmediatly schema key.
	AllowDeleteImmediatly = "allow-delete-immediatly"

	// ShowComputerIcon schema key.
	ShowComputerIcon string = "show-computer-icon"

	// ShowTrashIcon schema key.
	ShowTrashIcon = "show-trash-icon"

	// StickupGrid schema key.
	StickupGrid = "stickup-grid"

	// AutoArrangement schema key.
	AutoArrangement = "auto-arrangement"

	// ShowTrashedItemCount schema key.
	ShowTrashedItemCount = "show-trashed-item-number"

	// SortOrder schema key.
	SortOrder = "sort-order"

	// IconDefaultSize schema key
	IconDefaultSize = "icon-default-size"

	// IconZoomLevel schema key
	IconZoomLevel = "icon-zoom-level"

	// ThumbnailSizeLimitation schema key
	ThumbnailSizeLimitation = "thumbnail-size-limitation"

	// ThumbnailSizeUnit schema key
	ThumbnailSizeUnit = "thumbnail-size-unit"
)

const (
	// ActivationPolicyAsk indicates ask for behaviours when activation.
	ActivationPolicyAsk string = "ask"

	// ActivationPolicyLaunch indicates launch files when activation.
	ActivationPolicyLaunch string = "launch"

	// ActivationPolicyDisplay indicates display files when activation.
	ActivationPolicyDisplay string = "display"
)

// sortPoliciesName is a map to sort policies and the display name.
var sortPoliciesName = map[string]string{
	"name":      Tr("Name"),
	"size":      Tr("Size"),
	"filetype":  Tr("Filetype"),
	"mtime":     Tr("Modified time"),
	"atime":     Tr("Accessed time"),
	"open-with": Tr("Open with"),
	"tag-info":  Tr("Tag info"),
	"tag-color": Tr("Tag Color"),
}

const (
	SizeUnitByte int = iota
	SizeUnitKiB
	SizeUnitMB
	SizeUnitGB
	SizeUnitTB
	SizeUnitPB
)

func toBytes(size uint64, unit int) uint64 {
	switch unit {
	case SizeUnitByte:
		return size
	case SizeUnitKiB:
		return size * 1024
	case SizeUnitMB:
		return toBytes(size*1024, SizeUnitKiB)
	case SizeUnitGB:
		return toBytes(size*1024, SizeUnitMB)
	case SizeUnitTB:
		return toBytes(size*1024, SizeUnitGB)
	case SizeUnitPB:
		return toBytes(size*1024, SizeUnitTB)
	}
	panic("toBytes: invalid size unit, shouldn't reach here.")
}

// Settings is settings used by desktop.
type Settings struct {
	// filemanagerPreferences is filemanager's general preferences
	filemanagerPreferences *gio.Settings

	// desktopPreferences is desktop specific settings.
	desktopPreferences *gio.Settings

	// iconSize is the real icon size which equals to default icon size muliplies zoom level.
	iconSize int

	// thumbnailSizeLimitation is size limitation in bytes.
	thumbnailSizeLimitation uint64

	// signals
	IconZoomLevelChanged           func(int32)
	ShowTrashIconChanged           func(bool)
	ShowComputerIconChanged        func(bool)
	StickupGridChanged             func(bool)
	AutoArrangementChanged         func(bool)
	ConfirmEmptyTrashChanged       func(bool)
	ActivationPolicyChanged        func(string)
	ClickPolicyChanged             func(string)
	ShowThumbnailChanged           func(string)
	ShowHiddenFilesChanged         func(bool)
	ShowExtensionNameChanged       func(bool)
	LabelPositionChanged           func(string)
	AllowDeleteImmediatlyChanged   func(bool)
	ThumbnailSizeLimitationChanged func(uint64)
}

// GetDBusInfo returns dbus info for Settings.
func (s *Settings) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Desktop.Settings",
		ObjectPath: "/com/deepin/daemon/Desktop/Settings",
		Interface:  "com.deepin.daemon.Desktop.Settings",
	}
}

// NewSettings creates a new settings.
func NewSettings() (*Settings, error) {
	s := new(Settings)
	err := initializer.NewInitializer().Init(func(interface{}) (interface{}, error) {
		return utils.CheckAndNewGSettings(FileManagerPerferenceSchemaID)
	}).Init(func(v interface{}) (interface{}, error) {
		s.filemanagerPreferences = v.(*gio.Settings)
		s.filemanagerPreferences.Connect("changed", func(_ *gio.Settings, key string) {
			switch key {
			case ConfirmEmptyTrash:
				s.emitConfirmEmptyTrashChanged(s.ConfirmEmptyTrashIsEnable())
			case ConfirmTrash:
			case ActivationPolicy:
				s.emitActivationPolicyChanged(s.ActivationPolicy())
			case ClickPolicy:
				s.emitClickPolicyChanged(s.ClickPolicy())
			case ShowThumbnail:
				s.emitShowThumbnailChanged(s.ShowThumbnail())
			case ShowHiddenFiles:
				s.emitShowHiddenFilesChanged(s.ShowHiddenFilesIsEnable())
			case ShowExtensionName:
				s.emitShowExtensionNameChanged(s.ShowExtensionNameIsEnable())
			case LabelPosition:
				s.emitLabelPositionChanged(s.LabelPosition())
			case AllowDeleteImmediatly:
				s.emitAllowDeleteImmediatlyChanged(s.AllowDeleteImmediatlyIsEnable())
			case ThumbnailSizeLimitation:
				fallthrough
			case ThumbnailSizeUnit:
				s.updateThumbnailSizeLimitation()
				s.emitThunbnailSizeLimitationChanged(s.thumbnailSizeLimitation)
			}
		})
		s.filemanagerPreferences.GetBoolean(ConfirmEmptyTrash) // enable connection.
		s.updateThumbnailSizeLimitation()
		return utils.CheckAndNewGSettings(FileManagerDesktopSchemaID)
	}).Init(func(v interface{}) (interface{}, error) {
		s.desktopPreferences = v.(*gio.Settings)
		s.desktopPreferences.Connect("changed", func(_ *gio.Settings, key string) {
			switch key {
			case ShowComputerIcon:
				s.emitShowComputerIconChanged(s.ShowComputerIconIsEnable())
			case ShowTrashIcon:
				s.emitShowTrashIconChanged(s.ShowTrashIconIsEnable())
			case StickupGrid:
				s.emitStickupGridChanged(s.StickupGridIsEnable())
			// case ShowTrashedItemCount:
			case AutoArrangement:
				s.emitAutoArrangementChanged(s.AutoArrangement())
			case IconDefaultSize:
				s.updateIconSize()
			case IconZoomLevel:
				s.updateIconSize()
				s.emitIconZoomLevelChanged(s.IconZoomLevel())
			}
		})
		s.updateIconSize()
		return nil, nil
	}).GetError()

	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Settings) updateThumbnailSizeLimitation() {
	size := s.getThumbnailSizeLimitation()
	unit := int(s.getThumbnailSizeUnit())
	s.thumbnailSizeLimitation = toBytes(size, unit)
}

func (s *Settings) updateIconSize() {
	s.iconSize = int(s.desktopPreferences.GetEnum(IconDefaultSize) * s.desktopPreferences.GetInt(IconZoomLevel) / 100)
}

func (s *Settings) emitThunbnailSizeLimitationChanged(size uint64) error {
	return dbus.Emit(s, "ThumbnailSizeLimitationChanged", size)
}

func (s *Settings) emitIconZoomLevelChanged(level int32) error {
	return dbus.Emit(s, "IconZoomLevelChanged", level)
}

func (s *Settings) emitShowTrashIconChanged(enable bool) error {
	return dbus.Emit(s, "ShowTrashIconChanged", enable)
}

func (s *Settings) emitShowComputerIconChanged(enable bool) error {
	return dbus.Emit(s, "ShowComputerIconChanged", enable)
}

func (s *Settings) emitConfirmEmptyTrashChanged(enable bool) error {
	return dbus.Emit(s, "ConfirmEmptyTrashChanged", enable)
}

func (s *Settings) emitStickupGridChanged(enable bool) error {
	return dbus.Emit(s, "StickupGridChanged", enable)
}

func (s *Settings) emitAutoArrangementChanged(enable bool) error {
	return dbus.Emit(s, "AutoArrangementChanged", enable)
}

func (s *Settings) emitActivationPolicyChanged(activationPolicy string) error {
	return dbus.Emit(s, "ActivationPolicyChanged", activationPolicy)
}

func (s *Settings) emitClickPolicyChanged(clickPolicy string) error {
	return dbus.Emit(s, "ClickPolicyChanged", clickPolicy)
}

func (s *Settings) emitShowThumbnailChanged(showPolicy string) error {
	return dbus.Emit(s, "ShowThumbnailChanged", showPolicy)
}

func (s *Settings) emitShowHiddenFilesChanged(enable bool) error {
	return dbus.Emit(s, "ShowHiddenFilesChanged", enable)
}

func (s *Settings) emitShowExtensionNameChanged(enable bool) error {
	return dbus.Emit(s, "ShowExtensionNameChanged", enable)
}

func (s *Settings) emitLabelPositionChanged(position string) error {
	return dbus.Emit(s, "LabelPositionChanged", position)
}

func (s *Settings) emitAllowDeleteImmediatlyChanged(enable bool) error {
	return dbus.Emit(s, "AllowDeleteImmediatlyChanged", enable)
}

// ConfirmTrashIsEnable returns whether ConfirmTrash is enabled.
func (s *Settings) ConfirmTrashIsEnable() bool {
	return s.filemanagerPreferences.GetBoolean(ConfirmTrash)
}

// ShowTrashedItemCountIsEnable returns whether ShowTrashedItemCount is enabled.
func (s *Settings) ShowTrashedItemCountIsEnable() bool {
	return s.desktopPreferences.GetBoolean(ShowTrashedItemCount)
}

// ConfirmEmptyTrashIsEnable returns whether ConfirmEmptyTrash is enabled.
func (s *Settings) ConfirmEmptyTrashIsEnable() bool {
	return s.desktopPreferences.GetBoolean(ConfirmEmptyTrash)
}

func (s *Settings) getSortPolicies() []string {
	variantValue := s.filemanagerPreferences.GetRange(SortOrder)
	_, policies := variantValue.GetChildValue(1).GetVariant().GetStrv()
	return policies
}

// ShowComputerIconIsEnable returns whether ShowComputerIcon is enabled.
func (s *Settings) ShowComputerIconIsEnable() bool {
	return s.desktopPreferences.GetBoolean(ShowComputerIcon)
}

// EnableShowComputerIcon enables or disables ShowComputerIcon.
func (s *Settings) EnableShowComputerIcon(enable bool) {
	s.desktopPreferences.SetBoolean(ShowComputerIcon, enable)
}

// ShowTrashIconIsEnable returns whether ShowTrashIcon is enabled.
func (s *Settings) ShowTrashIconIsEnable() bool {
	return s.desktopPreferences.GetBoolean(ShowTrashIcon)
}

// EnableShowTrashIcon enables or disables ShowTrashIcon.
func (s *Settings) EnableShowTrashIcon(enable bool) {
	s.desktopPreferences.SetBoolean(ShowTrashIcon, enable)
}

// StickupGridIsEnable returns whether StickupGrid is enabled.
func (s *Settings) StickupGridIsEnable() bool {
	return s.desktopPreferences.GetBoolean(StickupGrid)
}

// EnableStickupGrid enables or disables StickupGrid.
func (s *Settings) EnableStickupGrid(enable bool) {
	s.desktopPreferences.SetBoolean(StickupGrid, enable)
}

func (s *Settings) AutoArrangement() bool {
	return s.desktopPreferences.GetBoolean(AutoArrangement)
}

func (s *Settings) EnableAutoArrangement(enable bool) {
	s.desktopPreferences.SetBoolean(AutoArrangement, enable)
}

// ActivationPolicy returns activation policy.
func (s *Settings) ActivationPolicy() string {
	return s.filemanagerPreferences.GetString(ActivationPolicy)
}

// ClickPolicy returns click policy.
func (s *Settings) ClickPolicy() string {
	return s.filemanagerPreferences.GetString(ClickPolicy)
}

// ShowThumbnail returns show thumbnail policy.
func (s *Settings) ShowThumbnail() string {
	return s.filemanagerPreferences.GetString(ShowThumbnail)
}

// ShowHiddenFilesIsEnable returns whether ShowHiddenFiles is enabled.
func (s *Settings) ShowHiddenFilesIsEnable() bool {
	return s.filemanagerPreferences.GetBoolean(ShowHiddenFiles)
}

// ShowExtensionNameIsEnable returns whether ShowExtensionName is enabled.
func (s *Settings) ShowExtensionNameIsEnable() bool {
	return s.filemanagerPreferences.GetBoolean(ShowExtensionName)
}

// LabelPosition returns the label position of name.
func (s *Settings) LabelPosition() string {
	return s.filemanagerPreferences.GetString(LabelPosition)
}

// AllowDeleteImmediatlyIsEnable returns whether AllowDeleteImmediatly is enabled.
func (s *Settings) AllowDeleteImmediatlyIsEnable() bool {
	return s.filemanagerPreferences.GetBoolean(AllowDeleteImmediatly)
}

// IconDefaultSize returns the default icon size.
func (s *Settings) IconDefaultSize() int32 {
	return s.desktopPreferences.GetInt(IconDefaultSize)
}

// IconZoomLevel returns the zoom level of icons.
func (s *Settings) IconZoomLevel() int32 {
	return s.desktopPreferences.GetInt(IconZoomLevel)
}

// SetIconZoomLevel will change the zoom level of icons.
func (s *Settings) SetIconZoomLevel(zoomLevel int32) bool {
	return s.desktopPreferences.SetInt(IconZoomLevel, zoomLevel)
}

func (s *Settings) getThumbnailSizeUnit() int32 {
	return s.filemanagerPreferences.GetEnum(ThumbnailSizeUnit)
}

func (s *Settings) getThumbnailSizeLimitation() uint64 {
	size := s.filemanagerPreferences.GetValue(ThumbnailSizeLimitation)
	return size.GetUint64()
}

func (s *Settings) ThumbnailSizeLimitation() uint64 {
	return s.thumbnailSizeLimitation
}
