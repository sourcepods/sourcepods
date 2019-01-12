part of swagger.api;

class Repository {
  
  String id = null;
  

  String name = null;
  

  String description = null;
  

  String website = null;
  

  String defaultBranch = null;
  

  DateTime createdAt = null;
  

  DateTime updatedAt = null;
  

  User owner = null;
  
  Repository();

  @override
  String toString() {
    return 'Repository[id=$id, name=$name, description=$description, website=$website, defaultBranch=$defaultBranch, createdAt=$createdAt, updatedAt=$updatedAt, owner=$owner, ]';
  }

  Repository.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    id =
        json['id']
    ;
    name =
        json['name']
    ;
    description =
        json['description']
    ;
    website =
        json['website']
    ;
    defaultBranch =
        json['defaultBranch']
    ;
    createdAt = json['createdAt'] == null ? null : DateTime.parse(json['createdAt']);
    updatedAt = json['updatedAt'] == null ? null : DateTime.parse(json['updatedAt']);
    owner =
      
      
      new User.fromJson(json['owner'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'website': website,
      'defaultBranch': defaultBranch,
      'createdAt': createdAt == null ? '' : createdAt.toUtc().toIso8601String(),
      'updatedAt': updatedAt == null ? '' : updatedAt.toUtc().toIso8601String(),
      'owner': owner
     };
  }

  static List<Repository> listFromJson(List<dynamic> json) {
    return json == null ? new List<Repository>() : json.map((value) => new Repository.fromJson(value)).toList();
  }

  static Map<String, Repository> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Repository>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Repository.fromJson(value));
    }
    return map;
  }
}

