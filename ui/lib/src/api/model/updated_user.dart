part of swagger.api;

class UpdatedUser {
  
  String name = null;
  
  UpdatedUser();

  @override
  String toString() {
    return 'UpdatedUser[name=$name, ]';
  }

  UpdatedUser.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name =
        json['name']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name
     };
  }

  static List<UpdatedUser> listFromJson(List<dynamic> json) {
    return json == null ? new List<UpdatedUser>() : json.map((value) => new UpdatedUser.fromJson(value)).toList();
  }

  static Map<String, UpdatedUser> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, UpdatedUser>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new UpdatedUser.fromJson(value));
    }
    return map;
  }
}

