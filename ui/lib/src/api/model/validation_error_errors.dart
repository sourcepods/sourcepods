part of swagger.api;

class ValidationErrorErrors {
  
  String field = null;
  

  String message = null;
  
  ValidationErrorErrors();

  @override
  String toString() {
    return 'ValidationErrorErrors[field=$field, message=$message, ]';
  }

  ValidationErrorErrors.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    field =
        json['field']
    ;
    message =
        json['message']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'field': field,
      'message': message
     };
  }

  static List<ValidationErrorErrors> listFromJson(List<dynamic> json) {
    return json == null ? new List<ValidationErrorErrors>() : json.map((value) => new ValidationErrorErrors.fromJson(value)).toList();
  }

  static Map<String, ValidationErrorErrors> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, ValidationErrorErrors>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new ValidationErrorErrors.fromJson(value));
    }
    return map;
  }
}

