import 'dart:async';
import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/routes.dart';
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/repository/repository_service.dart';

@Component(
  selector: 'gitpods-repository-create',
  templateUrl: 'repository_create_component.html',
  directives: const [coreDirectives, formDirectives],
  providers: const [RepositoryService],
)
class RepositoryCreateComponent {
  RepositoryCreateComponent(this._router, this._repositoryService);

  final Router _router;
  final RepositoryService _repositoryService;

  Repository repository = new Repository();
  bool loading;
  String error = '';

  void submit(Event event) {
    event.preventDefault();
    this.loading = true;
    this.error = '';

    _repositoryService
        .create(this.repository)
        .then((Repository repository) {
          _router.navigate(RoutePaths.repository.toUrl(parameters: {
            'owner': repository.owner.username,
            'name': repository.name,
          }));
        })
        .catchError((e) => this.error = e.toString())
        .whenComplete(() => loading = false);
  }
}
