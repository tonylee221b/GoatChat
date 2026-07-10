import 'package:flutter/material.dart';

import '../../../theme/app_theme.dart';
import '../../chat/presentation/chat_home_page.dart';

class NicknameEntryPage extends StatefulWidget {
  const NicknameEntryPage({super.key});

  @override
  State<NicknameEntryPage> createState() => _NicknameEntryPageState();
}

class _NicknameEntryPageState extends State<NicknameEntryPage> {
  final TextEditingController _controller = TextEditingController();

  bool get _canContinue => _controller.text.trim().length >= 2;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _handleNicknameChanged(String _) {
    setState(() {});
  }

  void _submit() {
    final nickname = _controller.text.trim();
    if (nickname.length < 2) {
      return;
    }

    Navigator.of(context).pushReplacement(
      MaterialPageRoute<void>(
        builder: (_) => ChatHomePage(nickname: nickname),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(20, 24, 20, 20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 32),
              const Text(
                'Before you start',
                style: TextStyle(
                  color: AppColors.mutedText,
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                ),
              ),
              const SizedBox(height: 10),
              const Text(
                'Enter your nickname',
                style: TextStyle(
                  color: AppColors.textPrimary,
                  fontSize: 32,
                  fontWeight: FontWeight.w800,
                ),
              ),
              const SizedBox(height: 12),
              Text(
                'This name will be shown in chat rooms.',
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: AppColors.mutedText,
                  height: 1.5,
                ),
              ),
              const SizedBox(height: 28),
              TextField(
                controller: _controller,
                autofocus: true,
                textInputAction: TextInputAction.done,
                onChanged: _handleNicknameChanged,
                onSubmitted: (_) => _submit(),
                decoration: InputDecoration(
                  hintText: 'Nickname',
                  hintStyle: const TextStyle(color: AppColors.mutedText),
                  filled: true,
                  fillColor: Colors.white,
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 20,
                    vertical: 18,
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(20),
                    borderSide: const BorderSide(
                      color: AppColors.border,
                      width: 1.2,
                    ),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(20),
                    borderSide: const BorderSide(
                      color: AppColors.border,
                      width: 1.2,
                    ),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(20),
                    borderSide: const BorderSide(
                      color: AppColors.ink,
                      width: 1.4,
                    ),
                  ),
                ),
              ),
              const Spacer(),
              SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton(
                  onPressed: _canContinue ? _submit : null,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.ink,
                    disabledBackgroundColor: AppColors.tabTrack,
                    foregroundColor: Colors.white,
                    disabledForegroundColor: AppColors.tabText,
                    elevation: 0,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(18),
                    ),
                    textStyle: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                  child: const Text('Continue'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
