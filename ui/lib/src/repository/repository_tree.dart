class RepositoryTree {
  String mode;
  String type;
  String object;
  String file;

  RepositoryTree({
    this.mode,
    this.type,
    this.object,
    this.file,
  });

  factory RepositoryTree.fromJSON(Map<String, dynamic> data) {
    return new RepositoryTree(
      mode: data['mode'],
      type: data['type'],
      object: data['object'],
      file: data['file'],
    );
  }
}
