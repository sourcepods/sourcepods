part of swagger.api;

class NewRepository {
  
  String name = null;
  

  String description = null;
  

  String website = null;
  
  NewRepository();

  @override
  String toString() {
    return 'NewRepository[name=$name, description=$description, website=$website, ]';
  }

  NewRepository.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name =
        json['name']
    ;
    description =
        json['description']
    ;
    website =
        json['website']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'description': description,
      'website': website
     };
  }

  static List<NewRepository> listFromJson(List<dynamic> json) {
    return json == null ? new List<NewRepository>() : json.map((value) => new NewRepository.fromJson(value)).toList();
  }

  static Map<String, NewRepository> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, NewRepository>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new NewRepository.fromJson(value));
    }
    return map;
  }
}

