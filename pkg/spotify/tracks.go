package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func (c Client) getNewTracksOfArtist(ctx context.Context, artist string) ([]spotify.SimpleTrack, error) {
	query := fmt.Sprintf("artist:%v tag:new", artist)
	res, err := c.spotify.Search(ctx, query, spotify.SearchTypeAlbum)
	if err != nil {
		return nil, fmt.Errorf("searching new album of %v: %w", artist, err)
	}
	var newTracks []spotify.SimpleTrack
	albums := res.Albums
	if albums == nil {
		return newTracks, nil
	}
	for _, album := range albums.Albums {
		tracks, err := c.getTracksOfAlbum(ctx, album.ID)
		if err != nil {
			return nil, fmt.Errorf("getting tracks of  album %v: %w", album.Name, err)
		}
		newTracks = append(newTracks, tracks...)
	}
	return newTracks, nil
}

func (c Client) getTracksOfAlbum(ctx context.Context, albumID spotify.ID) ([]spotify.SimpleTrack, error) {
	page, err := c.spotify.GetAlbumTracks(ctx, albumID)
	if err != nil {
		return nil, err
	}
	return page.Tracks, nil
}
