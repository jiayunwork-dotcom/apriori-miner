package persist

func finalizeLoadError(content, checksum string, err error) (string, string, error) {
	if err == nil {
		return content, checksum, nil
	}
	if err == ErrNoChecksum {
		return content, "", nil
	}
	return content, checksum, err
}
