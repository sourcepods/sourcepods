part of swagger.api;

class Commit {
  
  String hash = null;
  

  String tree = null;
  

  String parent = null;
  

  String message = null;
  

  String author = null;
  

  String authorEmail = null;
  

  String committer = null;
  

  String committerEmail = null;
  
  Commit();

  @override
  String toString() {
    return 'Commit[hash=$hash, tree=$tree, parent=$parent, message=$message, author=$author, authorEmail=$authorEmail, committer=$committer, committerEmail=$committerEmail, ]';
  }

  Commit.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    hash =
        json['hash']
    ;
    tree =
        json['tree']
    ;
    parent =
        json['parent']
    ;
    message =
        json['message']
    ;
    author =
        json['author']
    ;
    authorEmail =
        json['authorEmail']
    ;
    committer =
        json['committer']
    ;
    committerEmail =
        json['committerEmail']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'hash': hash,
      'tree': tree,
      'parent': parent,
      'message': message,
      'author': author,
      'authorEmail': authorEmail,
      'committer': committer,
      'committerEmail': committerEmail
     };
  }

  static List<Commit> listFromJson(List<dynamic> json) {
    return json == null ? new List<Commit>() : json.map((value) => new Commit.fromJson(value)).toList();
  }

  static Map<String, Commit> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Commit>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Commit.fromJson(value));
    }
    return map;
  }
}

