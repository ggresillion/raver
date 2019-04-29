import {Body, Controller, Get, Post} from '@nestjs/common';
import {Category} from './category.entity';
import {CategoryService} from './category.service';

@Controller('categories')
export class CategoryController {

  constructor(
    private readonly categoryService: CategoryService,
  ) {
  }

  @Get()
  public async getCategories(): Promise<Category[]> {
    return this.categoryService.findAllCategories();
  }

  @Post()
  public async createCategory(@Body()category: Category): Promise<any> {
    return this.categoryService.createCategory(category);
  }
}
