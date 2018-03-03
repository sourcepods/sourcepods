import 'package:angular/angular.dart';
import 'package:gitpods/src/loading_component.dart';

@Component(
  selector: 'repository-files',
  templateUrl: 'files_component.html',
  directives: const [
    COMMON_DIRECTIVES,
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
  String name;

  Branch(this.name);
}
