import 'package:sourcepods/src/api/api.dart' as api;

class TreeEntry {
  TreeEntry(this.mode, this.type, this.object, this.path);

  factory TreeEntry.fromAPI(api.TreeEntry te) {
    return TreeEntry(te.mode, te.type, te.object, te.path);
  }

  String mode;

  String type;

  String object;

  String path;
}
