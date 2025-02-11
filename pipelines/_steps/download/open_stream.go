package download

import (
	"errors"
	"io"

	"github.com/t2bot/matrix-media-repo/common/config"
	"github.com/t2bot/matrix-media-repo/common/rcontext"
	"github.com/t2bot/matrix-media-repo/database"
	"github.com/t2bot/matrix-media-repo/datastores"
	"github.com/t2bot/matrix-media-repo/redislib"
	"github.com/t2bot/matrix-media-repo/util/readers"
)

func OpenStream(ctx rcontext.RequestContext, media *database.Locatable) (io.ReadSeekCloser, error) {
	reader, ds, err := doOpenStream(ctx, media)
	if err != nil {
		return nil, err
	}
	if reader != nil {
		ctx.Log.Debugf("Got %s from cache", media.Sha256Hash)
		return readers.NopSeekCloser(reader), nil
	}

	return datastores.Download(ctx, ds, media.Location)
}

func OpenOrRedirect(ctx rcontext.RequestContext, media *database.Locatable) (io.ReadSeekCloser, error) {
	reader, ds, err := doOpenStream(ctx, media)
	if err != nil {
		return nil, err
	}
	if reader != nil {
		ctx.Log.Debugf("Got %s from cache", media.Sha256Hash)
		return readers.NopSeekCloser(reader), nil
	}

	return datastores.DownloadOrRedirect(ctx, ds, media.Location)
}

func doOpenStream(ctx rcontext.RequestContext, media *database.Locatable) (io.ReadSeekCloser, config.DatastoreConfig, error) {
	reader, err := redislib.TryGetMedia(ctx, media.Sha256Hash)
	if err != nil || reader != nil {
		ctx.Log.Debugf("Got %s from cache", media.Sha256Hash)
		return readers.NopSeekCloser(reader), config.DatastoreConfig{}, err
	}

	ds, ok := datastores.Get(ctx, media.DatastoreId)
	if !ok {
		return nil, ds, errors.New("unable to locate datastore for media")
	}

	return nil, ds, nil
}
