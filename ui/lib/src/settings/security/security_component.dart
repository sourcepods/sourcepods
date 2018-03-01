import 'package:angular/angular.dart';

@Component(
  selector: 'security',
  templateUrl: 'security_component.html',
)
class SecurityComponent implements OnInit {
  @override
  void ngOnInit() {
    print('Security');
  }
}
