package models

import "net/url"

/* //////////////////////////////
Struct and related methods to code logic
////////////////////////////// */

// Record Type used for code representation of a record
type Record struct {
	Target *url.URL // Link to be redirected to
	Id     string   // Short form URL
	Uid    string   // Record ID in database
}

/* //////////////////////////////
Generators
////////////////////////////// */

func CreateRecord(url *url.URL, id string, uid string) (*Record, error) {
	return &Record{
		Target: url,
		Id:     id,
		Uid:    uid,
	}, nil
}

func CreateRecordFromString(target string, id string, uid string) (*Record, error) {

	parsedUrl, urlErr := url.Parse(target)
	if urlErr != nil {
		return &Record{}, urlErr
	}

	return &Record{
		Target: parsedUrl,
		Id:     id,
		Uid:    uid,
	}, nil
}
