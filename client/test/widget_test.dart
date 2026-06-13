import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:goat_chat/src/app.dart';
import 'package:goat_chat/src/features/chat/presentation/chat_home_page.dart';

void main() {
  testWidgets('renders the GoatChat app shell without sample chat data', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const MyApp());

    expect(find.text('GoatChat'), findsOneWidget);
    expect(find.text('All Chats'), findsOneWidget);
    expect(find.text('Groups'), findsOneWidget);
    expect(find.text('Contacts'), findsOneWidget);
    expect(find.byType(ListView), findsNothing);
    expect(find.text('Larry Machigo'), findsNothing);
    expect(find.text('Natalie Nora'), findsNothing);

    expect(find.byIcon(Icons.call), findsNothing);
    expect(find.byIcon(Icons.videocam), findsNothing);
  });

  testWidgets('uses a modern monochrome palette', (WidgetTester tester) async {
    await tester.pumpWidget(const MyApp());

    final context = tester.element(find.byType(ChatHomePage));
    expect(Theme.of(context).colorScheme.primary, const Color(0xFF111111));
    expect(Theme.of(context).scaffoldBackgroundColor, const Color(0xFFF7F7F5));

    final fab = tester.widget<FloatingActionButton>(
      find.byType(FloatingActionButton),
    );
    expect(fab.backgroundColor, const Color(0xFF111111));
  });
}
