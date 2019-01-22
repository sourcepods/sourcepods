import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/routes.dart';
import 'package:sourcepods/src/loading_button_component.dart';
import 'package:sourcepods/src/repository/repository.dart';
import 'package:sourcepods/src/repository/repository_service.dart';

@Component(
  selector: 'gitpods-repository-create',
  templateUrl: 'repository_create_component.html',
  directives: const [coreDirectives, formDirectives, LoadingButtonComponent],
  providers: const [RepositoryService],
)
class RepositoryCreateComponent {
  RepositoryCreateComponent(this._router, this._repositoryService);

  final Router _router;
  final RepositoryService _repositoryService;

  bool loading;
  String error = '';

  String name;
  String nameError;
  String description;
  String descriptionError;
  String website;
  String websiteError;

  void submit(Event event) async {
    event.preventDefault();
    this.loading = true;
    this.error = '';
    this.nameError = null;
    this.descriptionError = null;
    this.websiteError = null;

    try {
      Repository r = await _repositoryService.create(
          this.name, this.description, this.website);
      _router.navigate(RoutePaths.repository.toUrl(parameters: {
        'owner': r.owner.username,
        'name': r.name,
      }));
    } on ValidationException catch (e) {
      _handleValidationException(e);
    } catch (e) {
      this.error = e.toString();
    } finally {
      this.loading = false;
    }
  }

  void _handleValidationException(ValidationException e) {
    this.error = '';

    e.errors.forEach((field, message) {
      switch (field) {
        case 'name':
          this.nameError = message;
          break;
        case 'description':
          this.descriptionError = message;
          break;
        case 'website':
          this.websiteError = message;
          break;
      }
    });
  }
}

class ValidationException implements Exception {
  ValidationException(this.message, this.errors);

  final String message;
  final Map<String, String> errors;

  factory ValidationException.fromJSON(Map<String, dynamic> data) {
    return ValidationException(
      data['message'],
      Map.fromIterable(
        data['errors'],
        key: (e) => e['field'],
        value: (e) => e['message'],
      ),
    );
  }
}
