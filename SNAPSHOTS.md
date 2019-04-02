# Persistent Snapshots

In addition to live streaming, `hkcam` lets you take snapshots.
Those snapshots are stored in the file system and can be loaded via custom characteristics.
[Persistent Snapshots](/SNAPSHOTS.md) are currently only supported by my [Home](https//hochgatterer.me/home) app.
If you are an app developer, you can implement the same functionality based on the following characteristic documentation.

## Characteristics

The following characteristics are used to take, get and delete snapshots.

- [TakeSnapshot](/take_snapshot.go) takes a snapshot.
- [Assets](/assets.go) returns an index of all snapshots as JSON.
- [GetAsset](/get_asset.go) returns raw bytes representing a snapshot JPEG image.
- [DeleteAssets](/delete_assets.go) deletes snapshots.

To take a snapshot, you should send `true` to the [TakeSnapshot](/take_snapshot.go) characteristic.

---

After that the [Assets](/assets.go) characteristic contains the snapshot in the returned JSON.
The content of the [Assets](/assets.go) characteritic might look like this.

```json
{
    "assets":[
        {
            "id": "1.jpg",
            "date": "2019-04-01T10:00:00+00:00"
        }
    ]
}
```

---

To get the bytes of the snapshot with id *1.jpg*, you should send the following JSON to the [GetAsset](/get_asset.go) characteristic.

```json
{
    "id":"1.jpg",
    "width":320,
    "height":240
}
```

After a successful write, you read the characteristic's value to get the JPEG bytes.

---

If you want to delete the snapshot, you send the following JSON to the [DeleteAssets](/delete_assets.go) characteristic.

```json
{
    "ids":[
        "1.jpg"
    ]
}
```