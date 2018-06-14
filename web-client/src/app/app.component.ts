import {Component, OnInit, ViewChild} from '@angular/core';
import {SongService} from './services/song/song.service';
import Category from './model/Category';
import {MatIconRegistry, MatSnackBar} from "@angular/material";
import {DomSanitizer} from "@angular/platform-browser";
import {forkJoin} from "rxjs/index";
import {BotService} from "./services/bot/bot.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {

  @ViewChild('file') file;
  public files: Set<File> = new Set();
  public categories: Category[] = [];
  public selectedCategory: string;
  public progress;
  public guilds = [];
  public selectedGuildId = 0;

  constructor(private songService: SongService,
              private matIconRegistry: MatIconRegistry,
              private domSanitizer: DomSanitizer,
              private snackBar: MatSnackBar,
              private botService: BotService) {
    matIconRegistry.addSvgIcon('songs', this.domSanitizer.bypassSecurityTrustResourceUrl('../assets/music-player.svg'));
    matIconRegistry.addSvgIcon('youtube', this.domSanitizer.bypassSecurityTrustResourceUrl('../assets/youtube.svg'));
    matIconRegistry.addSvgIcon('settings', this.domSanitizer.bypassSecurityTrustResourceUrl('../assets/settings.svg'));
    matIconRegistry.addSvgIcon('add', this.domSanitizer.bypassSecurityTrustResourceUrl('../assets/add.svg'));
  }

  ngOnInit() {
    this.refresh();
  }

  refresh() {
    this.songService.getSongs().subscribe((categories) => this.categories = categories);
    this.botService.getGuilds().subscribe((guilds) => this.guilds = guilds);
  }

  playSong(song: string, category: string) {
    this.songService.playSong(song, category).subscribe(() => {
    });
  }

  addSong(category: string) {
    this.file.nativeElement.click();
    this.selectedCategory = category;
  }

  onFilesAdded() {
    const files: { [key: string]: File } = this.file.nativeElement.files;
    for (const key in files) {
      if (!isNaN(parseInt(key, 10))) {
        this.files.add(files[key]);
      }
    }
    this.progress = this.songService.uploadSong(this.selectedCategory, this.files);
    let allProgressObservables = [];
    for (let key in this.progress) {
      allProgressObservables.push(this.progress[key].progress);
    }
    forkJoin(allProgressObservables).subscribe(() => {
      this.files = undefined;
      this.refresh();
      this.snackBar.open('Upload Completed !', '', {
        duration: 3000
      });
    });
  }

  selectGuild(guild) {
    this.selectedGuildId = guild.id;
  }

  getGuildIconClass(guild){
    const icon = 'round guild-icon';
    const selectedIcon = 'selected-guild round guild-icon';
    if(!this.selectedGuildId){
      return icon;
    }
    if(guild.id === this.selectedGuildId){
      return selectedIcon;
    } else {
      return icon;
    }
  }

  createCategory() {

  }
}
