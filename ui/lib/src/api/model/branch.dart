part of swagger.api;

class Branch {
  
  String name = null;
  

  String sha1 = null;
  

  String type = null;
  
  Branch();

  @override
  String toString() {
    return 'Branch[name=$name, sha1=$sha1, type=$type, ]';
  }

  Branch.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name =
        json['name']
    ;
    sha1 =
        json['sha1']
    ;
    type =
        json['type']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'sha1': sha1,
      'type': type
     };
  }

  static List<Branch> listFromJson(List<dynamic> json) {
    return json == null ? new List<Branch>() : json.map((value) => new Branch.fromJson(value)).toList();
  }

  static Map<String, Branch> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Branch>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Branch.fromJson(value));
    }
    return map;
  }
}

