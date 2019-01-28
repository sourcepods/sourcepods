import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/src/settings/routes.dart';

@Component(
  selector: 'settings',
  templateUrl: 'settings_component.html',
  directives: const [routerDirectives],
  exports: [Routes],
)
class SettingsComponent implements OnActivate {
  String tab;

  @override
  void onActivate(RouterState previous, RouterState current) {
    _setTab(current.path);
  }

  void _setTab(String path) {
    String tab = 'profile'; // default tab
    List<String> elements = path.split('/');

    // if tab other than default selected use this one
    if (elements.length > 2) {
      tab = elements.elementAt(2);
    }

    this.tab = tab;
  }

  String profileUrl() => RoutesPaths.profile.toUrl();

  String accountUrl() => RoutesPaths.account.toUrl();

  String securityUrl() => RoutesPaths.security.toUrl();

  bool tabActive(String tab) => this.tab == tab;
}
