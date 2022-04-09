package message

// FlagVisibleToCallerOnly is a flag that allows you to create messages visible only for the caller of the command (i.e. the user who triggered the command).
var FlagVisibleToCallerOnly uint64 = 1 << 6
