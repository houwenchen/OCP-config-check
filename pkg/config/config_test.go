package config

import "testing"

func TestModify_hosts(t *testing.T) {
	path := "./testdata/hosts"
	newPath := "./testdata/hosts_new"
	modify_hosts(path, newPath)
}

func TestWgetTar(t *testing.T) {
	download_scripts()
	tar_scripts()
}
