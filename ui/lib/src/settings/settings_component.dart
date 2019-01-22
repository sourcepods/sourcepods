import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/src/settings/routes.dart';

@Component(
  selector: 'settings',
  templateUrl: 'settings_component.html',
  directives: const [routerDirectives],
  exports: [Routes],
)
class SettingsComponent {
  String profileUrl() => RoutesPaths.profile.toUrl();

  String accountUrl() => RoutesPaths.account.toUrl();

  String securityUrl() => RoutesPaths.security.toUrl();
}
