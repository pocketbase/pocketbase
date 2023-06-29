package hooks

import (
	"context"
	"fmt"
	typesDocker "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/gorchitech/types"
	Gtype "github.com/pocketbase/pocketbase/gorchitech/types"
	"github.com/pocketbase/pocketbase/models"
	"io"
	"log"
	"os"
	"strings"
)

func SetupOnBeforeImagesHooks(app *pocketbase.PocketBase) {
	onBeforeCreateImage(app)
	onBeforeDeleteImage(app)
	onBeforeUpdateImage(app)
}

func onBeforeCreateImage(app *pocketbase.PocketBase) {
	app.OnModelBeforeCreate("image").Add(func(e *core.ModelEvent) error {
		//Create the docker containers for the service function
		gimage := types.GImageRecord{}
		err := gimage.FromRecordToImage(e.Model.(*models.Record))
		log.Println("on  SetupOnBeforeImagesHooks", gimage)

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		reader, err := cli.ImagePull(ctx, gimage.Name, typesDocker.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, reader)
		if err != nil {
			log.Fatalln("on  bytesToMB", err)
			panic(err)
		}
		defer reader.Close()

		// Get the pulled image ID
		imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), gimage.Name)
		if err != nil {
			panic(err)
		}
		gimage.ImageID = strings.Split(imageInspect.ID, ":")[1]
		gimage.Size = bytesToMB(imageInspect.Size)
		log.Println("on  bytesToMB", bytesToMB(imageInspect.Size))
		Gtype.FromImageToRecord(e.Model.(*models.Record), &gimage)
		log.Println("on  Images pull", imageInspect.Size)

		return nil
	})
}

func onBeforeDeleteImage(app *pocketbase.PocketBase) {
	app.OnModelBeforeDelete("image").Add(func(e *core.ModelEvent) error {
		log.Println("onBeforeDeleteImage", e.Dao, e.Model, e.Tags())
		gimage := types.GImageRecord{}
		err := gimage.FromRecordToImage(e.Model.(*models.Record))
		log.Println("on  SetupOnBeforeImagesHooks", gimage)
		if err != nil {
			panic(err)
		}
		log.Println("on  OnModelBeforeDelete")

		log.Println("on  Create CLI docker")
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		log.Println("on  Remove image")

		// Remove the image
		rsp, err := cli.ImageRemove(context.Background(), gimage.ImageID, typesDocker.ImageRemoveOptions{
			Force:         false, // Force removal even if the image is in use by containers
			PruneChildren: false, // Remove all images dependent on the specified image
		})
		if err != nil {
			return err
		}
		log.Println("on  Remove image", rsp)
		return nil
	})
}

func onBeforeUpdateImage(app *pocketbase.PocketBase) {
	app.OnModelBeforeUpdate("image").Add(func(e *core.ModelEvent) error {
		log.Println("onBeforeUpdateImage", e.Dao, e.Model, e.Tags())
		gimage := types.GImageRecord{}
		err := gimage.FromRecordToImage(e.Model.(*models.Record))
		log.Println("on  OnModelBeforeUpdate", gimage)
		if err != nil {
			panic(err)
		}
		log.Println("on  Create CLI docker")
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		if gimage.Deleted != "" {
			// Remove the image
			rsp, err := cli.ImageRemove(context.Background(), gimage.ImageID, typesDocker.ImageRemoveOptions{
				Force:         false, // Force removal even if the image is in use by containers
				PruneChildren: false, // Remove all images dependent on the specified image
			})
			if err != nil {
				return err
			}
			log.Println("on  Remove image", rsp)
			return nil
		}

		log.Println("on  Remove image")
		rsp, err := cli.ImageRemove(context.Background(), gimage.ImageID, typesDocker.ImageRemoveOptions{
			Force:         false, // Force removal even if the image is in use by containers
			PruneChildren: false, // Remove all images dependent on the specified image
		})
		if err != nil {
			return err
		}
		log.Println("on  Remove image", rsp)
		reader, err := cli.ImagePull(context.Background(), gimage.Name, typesDocker.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, reader)
		if err != nil {
			log.Fatalln("on  bytesToMB", err)
			panic(err)
		}
		defer reader.Close()

		// Get the pulled image ID
		imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), gimage.Name)
		if err != nil {
			panic(err)
		}
		gimage.ImageID = strings.Split(imageInspect.ID, ":")[1]
		gimage.Size = bytesToMB(imageInspect.Size)
		log.Println("on  bytesToMB", bytesToMB(imageInspect.Size))
		Gtype.FromImageToRecord(e.Model.(*models.Record), &gimage)
		log.Println("on  Images pull", imageInspect.Size)
		return nil
	})
}

func bytesToMB(bytes int64) string {
	mb := float64(bytes) / (1024 * 1024)
	log.Println("on  bytesToMB", mb)
	return fmt.Sprintf("%.2f MB", mb)
}
