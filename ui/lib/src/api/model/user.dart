part of swagger.api;

class User {
  
  String id = null;
  

  String email = null;
  

  String username = null;
  

  String name = null;
  

  DateTime createdAt = null;
  

  DateTime updatedAt = null;
  
  User();

  @override
  String toString() {
    return 'User[id=$id, email=$email, username=$username, name=$name, createdAt=$createdAt, updatedAt=$updatedAt, ]';
  }

  User.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    id =
        json['id']
    ;
    email =
        json['email']
    ;
    username =
        json['username']
    ;
    name =
        json['name']
    ;
    createdAt = json['createdAt'] == null ? null : DateTime.parse(json['createdAt']);
    updatedAt = json['updatedAt'] == null ? null : DateTime.parse(json['updatedAt']);
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'username': username,
      'name': name,
      'createdAt': createdAt == null ? '' : createdAt.toUtc().toIso8601String(),
      'updatedAt': updatedAt == null ? '' : updatedAt.toUtc().toIso8601String()
     };
  }

  static List<User> listFromJson(List<dynamic> json) {
    return json == null ? new List<User>() : json.map((value) => new User.fromJson(value)).toList();
  }

  static Map<String, User> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, User>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new User.fromJson(value));
    }
    return map;
  }
}

