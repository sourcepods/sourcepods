import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/src/loading_component.dart';
import 'package:sourcepods/src/repository/repository_service.dart';
import 'package:sourcepods/src/repository/tree.dart';

@Component(
  selector: 'repository-files',
  templateUrl: 'files_component.html',
  directives: const [
    coreDirectives,
    routerDirectives,
    formDirectives,
    LoadingComponent,
  ],
  providers: [
    ClassProvider(RepositoryService),
  ],
)
class FilesComponent implements OnActivate {
  FilesComponent(this._repositoryService);

  RepositoryService _repositoryService;

  String ownerName;
  String repositoryName;
  String defaultBranch = 'master'; // TODO: needs to be @Input()
  String currentBranch;
  String path;

  bool loading;
  List<String> branches;
  List<TreeEntry> files;
  List<TreeEntry> folders;

  @override
  void onActivate(RouterState previous, RouterState current) {
    this.ownerName = current.parameters['owner'];
    this.repositoryName = current.parameters['name'];
    this._fromPath(current.path);

    Future.wait([
      _repositoryService.getBranches(
        ownerName,
        repositoryName,
      ),
      _repositoryService.getTree(
        ownerName,
        repositoryName,
        currentBranch,
        path,
      ),
    ]).then((List responses) {
      _setBranches(responses[0]);
      _setTree(responses[1]);
    }).whenComplete(() => this.loading = false);
  }

  void _fromPath(String path) {
    List<String> elements = path.split('/');

    if (elements.length <= 4) {
      this.currentBranch = this.defaultBranch;
      this.path = '.';
      return;
    }

    // First elements are owner, name and tree

    currentBranch = elements.elementAt(4);
    this.path = elements.sublist(5).join('/');
  }

  void _setBranches(List<String> branches) {
    this.branches = branches;
  }

  void _setTree(List<TreeEntry> tree) {
    List<TreeEntry> files = [];
    List<TreeEntry> folders = [];

    tree.forEach((te) {
      if (te.type == 'tree') {
        folders.add(te);
      } else {
        files.add(te);
      }
    });

    files.sort((f1, f2) => f1.path.compareTo(f2.path));
    folders.sort((f1, f2) => f1.path.compareTo(f2.path));

    this.files = files;
    this.folders = folders;
  }

  void changeBranch(Event e) {
    print(e);
  }

  String changePath(String path) {
    return './$ownerName/$repositoryName/tree/$currentBranch/$path';
  }

  String openBlob(String path) {
    return './$ownerName/$repositoryName/blob/$currentBranch/$path';
  }
}
