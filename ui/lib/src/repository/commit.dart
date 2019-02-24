import 'package:sourcepods/src/api/api.dart' as api;

class Commit {
  Commit(this.hash, this.tree, this.message, this.author);

  factory Commit.fromAPI(api.Commit commit) {
    return Commit(commit.hash, commit.tree, commit.message, commit.author);
  }

  String hash;

  String tree;

  String message;

  String author;
}
