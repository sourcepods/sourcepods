class Repository {
  String id;
  String name;
  String description;
  String website;
  String default_branch;
  bool private;
  bool bare;
  DateTime created;
  DateTime updated;

  int stars;
  int forks;

  Repository({
    this.id,
    this.name,
    this.description,
    this.website,
    this.default_branch,
    this.private,
    this.bare,
    this.created,
    this.updated,
    this.stars,
    this.forks,
  });
}
