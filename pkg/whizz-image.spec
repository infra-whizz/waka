#
# spec file for whizz-image package
#

Name:           whizz-image
Version:        0.9
Release:        0
Summary:        Whizz Image Generator
License:        MIT
Group:          System/Tools
Url:            https://gitlab.com/infra-whizz/waka
Source:         %{name}-%{version}.tar.gz
Source1:        vendor.tar.gz

BuildRequires:  golang-packaging
BuildRequires:  golang(API) >= 1.13
Requires:       whizz-client
Requires:       whizz-ansible-modules

%description
Whizz-based image generator, which is using generic Ansible modules

%prep
%setup -q
%setup -q -T -D -a 1

%build
go build -x -mod=vendor -buildmode=pie -o %{name} ./cmd/*.go

%install
install -D -m 0755 %{name} %{buildroot}%{_bindir}/%{name}
mkdir -p %{buildroot}%{_sysconfdir}
install -m 0644 ./etc/waka.conf %{buildroot}%{_sysconfdir}/%{name}.conf

%files
%defattr(-,root,root)
%{_bindir}/%{name}
%dir %{_sysconfdir}
%config /etc/%{name}.conf

%changelog
