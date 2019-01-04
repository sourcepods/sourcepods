import 'package:angular/angular.dart';
import 'package:gitpods/src/loading_component.dart';

@Component(
  selector: 'repository-files',
  templateUrl: 'files_component.html',
  directives: const [
    coreDirectives,
    LoadingComponent,
  ],
)
class FilesComponent {
  String defaultBranch = 'master'; // TODO: needs to be @Input()

  bool loading;
  List<Branch> branches = [
    new Branch('master'),
    new Branch('develop'),
  ];
}

class Branch {
  Branch(this.name);

  String name;
}
