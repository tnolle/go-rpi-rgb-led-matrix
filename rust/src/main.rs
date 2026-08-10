use ratatui::widgets::{Block, Paragraph};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    ratatui::run(|terminal| {
        terminal.draw(|frame| {
            let block = Block::bordered().title("Welcome");
            let greeting = Paragraph::new("Hello, Ratatui! 🐭").centered().block(block);
            frame.render_widget(greeting, frame.area());
        })?;
        std::thread::sleep(std::time::Duration::from_secs(5));
        Ok(())
    })
}
