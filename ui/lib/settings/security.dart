import 'package:angular/angular.dart';

@Component(
  selector: 'security',
  templateUrl: 'security.html',
)
class SecurityComponent implements OnInit {
  @override
  void ngOnInit() {
    print('Security');
  }
}
