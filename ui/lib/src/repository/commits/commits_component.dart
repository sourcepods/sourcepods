// import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/src/loading_component.dart';
import 'package:sourcepods/src/repository/files/breadcrumb_component.dart';
import 'package:sourcepods/src/repository/repository_service.dart';
import 'package:sourcepods/src/repository/commit.dart';

@Component(
  selector: 'repository-commits',
  templateUrl: 'commits_component.html',
  directives: const [
    coreDirectives,
    routerDirectives,
    formDirectives,
    LoadingComponent,
    BreadcrumbComponent,
  ],
  providers: [
    ClassProvider(RepositoryService),
  ],
)
class CommitsComponent implements OnActivate {
  CommitsComponent(this._repositoryService);
  
  RepositoryService _repositoryService;

  String ownerName;
  String repositoryName;
  String defaultBranch = 'master'; // TODO: needs to be @Input
  
  bool loading;
  List<Commit> commits;

  @override
  void onActivate(RouterState previous, RouterState current) {
    this.ownerName =current.parameters['owner'];
    this.repositoryName =current.parameters['name'];

    Future.wait([
      _repositoryService.getCommits(
        ownerName,
        repositoryName,
        defaultBranch,
      ),
    ]).then((List responses) {
      _setCommits(responses[0]);
    }).whenComplete(() => this.loading = false);
  }

  void _setCommits(List<Commit> commits) {
    this.commits = commits;
  }
}
