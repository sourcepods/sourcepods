import 'dart:async';
import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/repository.dart';
import 'package:gitpods/repository_service.dart';

@Component(
  selector: 'gitpods-repository-create',
  templateUrl: 'repository_create_component.html',
  directives: const [COMMON_DIRECTIVES, formDirectives],
  providers: const [RepositoryService],
)
class RepositoryCreateComponent {
  final Router _router;
  final RepositoryService _repositoryService;

  RepositoryCreateComponent(this._router, this._repositoryService);

  Repository repository = new Repository();
  bool loading;
  String error = '';

  void submit(Event event) {
    event.preventDefault();
    this.loading = true;
    this.error = '';

    Future<Repository> resp = this._repositoryService.create(this.repository);

    resp.then((Repository repository) => this._router.navigate(['Repository', {'owner': repository.owner.username, 'name': repository.name}]))
        .catchError((e) => this.error = e.toString())
        .whenComplete(() => this.loading = false);
  }
}
