part of swagger.api;

class ValidationError {
  
  String message = null;
  

  List<ValidationErrorErrors> errors = [];
  
  ValidationError();

  @override
  String toString() {
    return 'ValidationError[message=$message, errors=$errors, ]';
  }

  ValidationError.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    message =
        json['message']
    ;
    errors =
      ValidationErrorErrors.listFromJson(json['errors'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'message': message,
      'errors': errors
     };
  }

  static List<ValidationError> listFromJson(List<dynamic> json) {
    return json == null ? new List<ValidationError>() : json.map((value) => new ValidationError.fromJson(value)).toList();
  }

  static Map<String, ValidationError> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, ValidationError>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new ValidationError.fromJson(value));
    }
    return map;
  }
}

