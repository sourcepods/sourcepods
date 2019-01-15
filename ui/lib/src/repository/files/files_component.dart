import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/loading_component.dart';
import 'package:gitpods/src/repository/repository_service.dart';

@Component(
  selector: 'repository-files',
  templateUrl: 'files_component.html',
  directives: const [
    coreDirectives,
    LoadingComponent,
  ],
  providers: [
    ClassProvider(RepositoryService),
  ],
)
class FilesComponent implements OnActivate {
  FilesComponent(this._repositoryService);

  RepositoryService _repositoryService;

  bool loading;
  String defaultBranch = 'master'; // TODO: needs to be @Input()
  List<String> branches;

  @override
  void onActivate(RouterState previous, RouterState current) {
    String ownerName = current.parameters['owner'];
    String name = current.parameters['name'];

    _repositoryService
        .getBranches(ownerName, name)
        .then((branches) => this.branches = branches)
        .whenComplete(() => this.loading = false);
  }
}
