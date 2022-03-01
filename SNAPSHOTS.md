# Persistent Snapshots

*Persistent Snapshots* is a way to take snapshots of the camera and store them on disk.
You can then access them via HomeKit.
*Persistent Snapshots* is not defined in the HAP but instead implemented by `hkcam` with custom characteristics.

*Persistent Snapshots* are supported by [Home+](https://hochgatterer.me/home+).

## Why?

Taking snapshots is an essential feature of any security camera.
For example you want to take a snapshot once motion is detected in room.
But there are currently no IP cameras which support that via HomeKit.

`hkcam` implements *Persistent Snapshots* with custom HomeKit characteristics.
This means you can use this feature in HomeKit scenes and automations.

## Custom Characteristics

The following characteristics are used to handle snapshots.

- [TakeSnapshot](/take_snapshot.go) takes a snapshot.
- [Assets](/assets.go) returns an index of all snapshots as JSON.
- [GetAsset](/get_asset.go) returns JPEG data representing a snapshot.
- [DeleteAssets](/delete_assets.go) deletes snapshots.

To take a snapshot, you should write `true` to the [TakeSnapshot](/take_snapshot.go) characteristic.

---

After that the [Assets](/assets.go) characteristic contains the snapshot in the returned JSON.
The value of the [Assets](/assets.go) characteristic might look like this.

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

To get the data of the snapshot with id `1.jpg`, you should send the following JSON to the [GetAsset](/get_asset.go) characteristic.

```json
{
    "id": "1.jpg",
    "width": 320,
    "height": 240
}
```

If you omit `width` or `height` (or set it to `0`), the image keeps the aspect ratio while resizing.

After a successful write, the [GetAsset](/get_asset.go) characteristic value contains the JPEG data.

---

If you want to delete the snapshot, you send the following JSON to the [DeleteAssets](/delete_assets.go) characteristic.

```json
{
    "ids": [
        "1.jpg"
    ]
}
```