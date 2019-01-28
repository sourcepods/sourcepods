part of swagger.api;

class TreeEntry {
  
  String mode = null;
  

  String type = null;
  

  String object = null;
  

  String path = null;
  
  TreeEntry();

  @override
  String toString() {
    return 'TreeEntry[mode=$mode, type=$type, object=$object, path=$path, ]';
  }

  TreeEntry.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    mode =
        json['mode']
    ;
    type =
        json['type']
    ;
    object =
        json['object']
    ;
    path =
        json['path']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'mode': mode,
      'type': type,
      'object': object,
      'path': path
     };
  }

  static List<TreeEntry> listFromJson(List<dynamic> json) {
    return json == null ? new List<TreeEntry>() : json.map((value) => new TreeEntry.fromJson(value)).toList();
  }

  static Map<String, TreeEntry> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, TreeEntry>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new TreeEntry.fromJson(value));
    }
    return map;
  }
}

